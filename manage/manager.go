package manage

import (
	"context"
	"errors"
	"fmt"

	"github.com/imrenagi/go-payment/config"
	"github.com/imrenagi/go-payment/datastore/inmemory"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/datastore"
	midgateway "github.com/imrenagi/go-payment/gateway/midtrans"
	xengateway "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/util/localconfig"

	"github.com/rs/zerolog"
)

// NewDefaultManager create use default properties for a manager instance
func NewDefaultManager(
	secret localconfig.PaymentSecret,
) *Manager {
	m := NewManager(secret)
	// default implementation
	m.paymentConfigRepository = inmemory.NewPaymentConfigRepository()

	return m
}

// NewManager creates a new payment manager
func NewManager(
	secret localconfig.PaymentSecret,
) *Manager {
	return &Manager{
		xenditGateway:   xengateway.NewGateway(secret.Xendit),
		midtransGateway: midgateway.NewGateway(secret.Midtrans),

		// invoiceRepository:        mysql.NewInvoiceRepository(db),
	}
}

type paymentConfigRepository interface {
	FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.Option) (config.FeeConfigReader, error)
	FindAll(ctx context.Context) (*config.PaymentConfig, error)
}

// Manager handle business logic related to payment gateway
type Manager struct {
	xenditGateway            *xengateway.Gateway
	midtransGateway          *midgateway.Gateway
	midTransactionRepository datastore.MidtransTransactionStatusRepository
	invoiceRepository        datastore.InvoiceRepository
	paymentConfigRepository  paymentConfigRepository
}

// MapMidtransTransactionStatusRepository mapping the midtrans transaction status repository
func (m *Manager) MapMidtransTransactionStatusRepository(repo datastore.MidtransTransactionStatusRepository) error {
	m.midTransactionRepository = repo
	return nil
}

// MustMidtransTransactionStatusRepository mandatory mapping the midtrans transaction status repo interface
func (m *Manager) MustMidtransTransactionStatusRepository(repo datastore.MidtransTransactionStatusRepository) {
	if repo == nil {
		panic(fmt.Errorf("midtrans transaction status repository can't be nil"))
	}
	m.midTransactionRepository = repo
}

// MustInvoiceRepository mapping the invoice repository
func (m *Manager) MustInvoiceRepository(repo datastore.InvoiceRepository) {
	if repo == nil {
		panic(fmt.Errorf("invoice repository can't be nil"))
	}
	m.invoiceRepository = repo
}

func (m Manager) charger(inv *invoice.Invoice) invoice.PaymentCharger {
	switch payment.NewGateway(inv.Payment.Gateway) {
	case payment.GatewayXendit:
		return &xenditCharger{
			XenditGateway: m.xenditGateway,
		}
	case payment.GatewayMidtrans:
		return &midtransCharger{
			MidtransGateway: m.midtransGateway,
		}
	default:
		panic("payment gateway is not found.")
	}
}

// GetPaymentMethods return the payment methods available in payment service
func (m *Manager) GetPaymentMethods(ctx context.Context, opts ...payment.Option) (*PaymentMethodList, error) {
	cfg, err := m.paymentConfigRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	options := payment.Options{}
	for _, o := range opts {
		o(&options)
	}

	return paymentMethodListFromConfig(cfg, options.Price), nil
}

// GetInvoice return invoice given its invoice number
func (m *Manager) GetInvoice(ctx context.Context, number string) (*invoice.Invoice, error) {
	return m.invoiceRepository.FindByNumber(ctx, number)
}

// GenerateInvoice generates new invoice
func (m *Manager) GenerateInvoice(ctx context.Context, gir *GenerateInvoiceRequest) (*invoice.Invoice, error) {

	var opts []payment.Option
	if gir.Payment.CreditCardDetail != nil {
		opts = append(opts, payment.WithCreditCard(
			gir.Payment.CreditCardDetail.Bank,
			gir.Payment.CreditCardDetail.Installment.Type,
			gir.Payment.CreditCardDetail.Installment.Term,
		))
	}

	paymentConfig, err := m.paymentConfigRepository.FindByPaymentType(ctx, gir.Payment.PaymentType, opts...)
	if err != nil {
		return nil, err
	}

	payment, err := invoice.NewPayment(paymentConfig, gir.Payment.PaymentType, gir.Payment.CreditCardDetail)
	if err != nil {
		return nil, err
	}

	inv := invoice.NewDefault()
	if err = inv.SetItem(ctx, *invoice.NewLineItem(
		gir.Item.Name,
		gir.Item.Category,
		gir.Item.MerchantName,
		gir.Item.Description,
		gir.Item.Price,
		gir.Item.Qty,
		gir.Item.Currency,
	)); err != nil {
		return nil, err
	}

	if err = inv.UpsertBillingAddress(gir.Customer.Name, gir.Customer.Email, gir.Customer.PhoneNumber); err != nil {
		return nil, err
	}

	if err = inv.UpdatePaymentMethod(ctx, payment, m.paymentConfigRepository, opts...); err != nil {
		return nil, err
	}

	if err = inv.Publish(ctx); err != nil {
		return nil, err
	}

	if err = inv.CreateChargeRequest(ctx, m.charger(inv)); err != nil {
		return nil, err
	}

	if err = m.invoiceRepository.Save(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

// PayInvoice pays an invoice. Invoice can only be paid if it is in right state
func (m *Manager) PayInvoice(ctx context.Context, pir *PayInvoiceRequest) (*invoice.Invoice, error) {

	log := zerolog.Ctx(ctx).With().
		Str("invoice_number", pir.InvoiceNumber).
		Logger()

	inv, err := m.invoiceRepository.FindByNumber(ctx, pir.InvoiceNumber)
	if errors.Is(err, payment.ErrNotFound) {
		return nil, nil
	}
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return nil, err
	}

	err = inv.Pay(ctx, pir.TransactionID)
	if err != nil {
		return nil, err
	}

	err = m.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("invoice paid")

	return inv, nil
}

// ProcessInvoice used if payment is initiated from user's end. It's either because they are using VA or any payment
// methods that requires payment action from the user after they choose a payment method/see payment instruction
func (m *Manager) ProcessInvoice(ctx context.Context, invoiceNumber string) (*invoice.Invoice, error) {

	inv, err := m.invoiceRepository.FindByNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, err
	}

	err = inv.Process(ctx)
	if err != nil {
		return nil, err
	}

	err = m.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

// FailInvoice fails an invoice if the payment is either failed or expired
func (m *Manager) FailInvoice(ctx context.Context, fir *FailInvoiceRequest) (*invoice.Invoice, error) {
	inv, err := m.invoiceRepository.FindByNumber(ctx, fir.InvoiceNumber)
	if errors.Is(err, payment.ErrNotFound) {
		return nil, nil
	}
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return nil, err
	}

	err = inv.Fail(ctx)
	if err != nil {
		return nil, err
	}

	err = m.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
}
