package manage

import (
	"context"
	"errors"
	"fmt"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/datastore"
	midgateway "github.com/imrenagi/go-payment/gateway/midtrans"
	xengateway "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/subscription"
	"github.com/imrenagi/go-payment/util/localconfig"

	"github.com/rs/zerolog"
)

// NewManager creates a new payment manager
func NewManager(
	secret localconfig.PaymentSecret,
) *Manager {
	return &Manager{
		xenditGateway:   xengateway.NewGateway(secret.Xendit),
		midtransGateway: midgateway.NewGateway(secret.Midtrans),
	}
}

// Manager handle business logic related to payment gateway
type Manager struct {
	xenditGateway            *xengateway.Gateway
	midtransGateway          *midgateway.Gateway
	midTransactionRepository datastore.MidtransTransactionStatusRepository
	invoiceRepository        datastore.InvoiceRepository
	subscriptionRepository   datastore.SubscriptionRepository
	paymentConfigRepository  datastore.PaymentConfigReader
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

// MustInvoiceRepository mandatory mapping the invoice repository
func (m *Manager) MustInvoiceRepository(repo datastore.InvoiceRepository) {
	if repo == nil {
		panic(fmt.Errorf("invoice repository can't be nil"))
	}
	m.invoiceRepository = repo
}

// MustSubscriptionRepository mandatory mapping the subscription repository
func (m *Manager) MustSubscriptionRepository(repo datastore.SubscriptionRepository) {
	if repo == nil {
		panic(fmt.Errorf("invoice repository can't be nil"))
	}
	m.subscriptionRepository = repo
}

// MustPaymentConfigReader mandatory mapping for payment config repository
func (m *Manager) MustPaymentConfigReader(repo datastore.PaymentConfigReader) {
	if repo == nil {
		panic(fmt.Errorf("invoice repository can't be nil"))
	}
	m.paymentConfigRepository = repo
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
	var items []invoice.LineItem
	for _, item := range gir.Items {
		i := invoice.NewLineItem(
			item.Name,
			item.Category,
			item.MerchantName,
			item.Description,
			item.Price,
			item.Qty,
			item.Currency,
		)
		items = append(items, *i)
	}
	if err = inv.SetItems(ctx, items); err != nil {
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

// CreateSubscription creates new subscription
func (m *Manager) CreateSubscription(ctx context.Context, csr *CreateSubscriptionRequest) (*subscription.Subscription, error) {
	s := csr.ToSubscription()

	if err := s.Start(ctx, m.subscriptionController(payment.GatewayXendit)); err != nil {
		return nil, err
	}

	if err := m.subscriptionRepository.Save(ctx, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (m Manager) subscriptionController(gateway payment.Gateway) subscription.Controller {
	return &xenditSubscriptionController{
		XenditGateway: m.xenditGateway,
	}
}
