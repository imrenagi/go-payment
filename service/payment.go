package service

import (
	"context"
	"errors"

	"github.com/imrenagi/go-payment/datastore/inmemory"
	"github.com/imrenagi/go-payment/util/localconfig"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/datastore"
	"github.com/imrenagi/go-payment/datastore/mysql"
	midgateway "github.com/imrenagi/go-payment/gateway/midtrans"
	xengateway "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

// NewService creates a new payment service
func NewService(
	db *gorm.DB,
	secret localconfig.PaymentSecret,
) *Service {
	return &Service{
		XenditGateway:            xengateway.NewGateway(secret.Xendit),
		MidtransGateway:          midgateway.NewGateway(secret.Midtrans),
		midCardTokenRepository:   &mysql.MidtransCardTokenRepository{DB: db},
		midTransactionRepository: &mysql.MidtransTransactionRepository{DB: db},
		invoiceRepository:        mysql.NewInvoiceRepository(db),
		paymentConfigRepository:  inmemory.NewPaymentConfigRepository(),
	}
}

// Service handle business logic related to payment gateway
type Service struct {
	XenditGateway            *xengateway.Gateway
	MidtransGateway          *midgateway.Gateway
	midCardTokenRepository   datastore.MidtransCardTokenRepository
	midTransactionRepository datastore.MidtransTransactionStatusRepository
	paymentConfigRepository  datastore.PaymentMethodRepository
	invoiceRepository        datastore.InvoiceRepository
}

func (p *Service) GetPaymentMethods(ctx context.Context, opts ...payment.PaymentOption) (*PaymentMethodList, error) {
	cfg, err := p.paymentConfigRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	options := payment.PaymentOptions{}
	for _, o := range opts {
		o(&options)
	}

	return paymentMethodListFromConfig(cfg, options.Price), nil
}

func (p *Service) GetInvoice(ctx context.Context, invoiceNumber string) (*invoice.Invoice, error) {
	return p.invoiceRepository.FindByNumber(ctx, invoiceNumber)
}

func (p *Service) Donate(ctx context.Context, cmd CreateDonationCommand) (*invoice.Invoice, error) {

	var opts []payment.PaymentOption
	if cmd.Payment.CreditCardDetail != nil {
		opts = append(opts, payment.WithCreditCard(
			cmd.Payment.CreditCardDetail.Bank,
			cmd.Payment.CreditCardDetail.Installment.Type,
			cmd.Payment.CreditCardDetail.Installment.Term,
		))
	}

	paymentConfig, err := p.paymentConfigRepository.FindByPaymentType(ctx, cmd.Payment.PaymentType, opts...)
	if err != nil {
		return nil, err
	}

	payment, err := invoice.NewPayment(paymentConfig, cmd.Payment.PaymentType, cmd.Payment.CreditCardDetail)
	if err != nil {
		return nil, err
	}

	inv := invoice.NewDefault(cmd.Message)

	if err = inv.SetItem(ctx, *invoice.NewLineItem(
		"Donasi",
		"PODCAST",
		"Ngobrolin Startup & Teknologi",
		cmd.Value,
		"IDR",
	)); err != nil {
		return nil, err
	}

	if err = inv.SetBillingAddress(cmd.DonaturName, cmd.Email, cmd.PhoneNumber); err != nil {
		return nil, err
	}

	if err = inv.SetPaymentMethod(payment); err != nil {
		return nil, err
	}

	if err = inv.UpdateFee(ctx, p.paymentConfigRepository, opts...); err != nil {
		return nil, err
	}

	if err = inv.Publish(ctx); err != nil {
		return nil, err
	}

	if err = inv.CreateChargeRequest(ctx, p.charger(inv)); err != nil {
		return nil, err
	}

	if err = p.invoiceRepository.Save(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

func (p Service) charger(inv *invoice.Invoice) invoice.PaymentCharger {
	switch inv.Payment.Gateway {
	case "xendit":
		return &xenditCharger{
			XenditGateway: p.XenditGateway,
		}
	default:
		return &midtransCharger{
			MidtransGateway: p.MidtransGateway,
		}
	}
}

func (p *Service) PayInvoice(ctx context.Context, invoiceNumber string, cmd PayInvoiceCommand) (*invoice.Invoice, error) {

	log := zerolog.Ctx(ctx).With().
		Str("invoice_number", invoiceNumber).
		Logger()

	inv, err := p.invoiceRepository.FindByNumber(ctx, invoiceNumber)
	if errors.Is(err, payment.ErrNotFound) {
		return nil, nil
	}
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return nil, err
	}

	err = inv.Pay(ctx, cmd.TransactionID)
	if err != nil {
		return nil, err
	}

	err = p.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("invoice paid")

	return inv, nil
}

func (p *Service) ProcessInvoice(ctx context.Context, invoiceNumber string) (*invoice.Invoice, error) {

	inv, err := p.invoiceRepository.FindByNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, err
	}

	err = inv.Process(ctx)
	if err != nil {
		return nil, err
	}

	err = p.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}

	return inv, nil

}

func (p *Service) FailInvoice(ctx context.Context, invoiceNumber string) (*invoice.Invoice, error) {
	inv, err := p.invoiceRepository.FindByNumber(ctx, invoiceNumber)
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

	err = p.invoiceRepository.Save(ctx, inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
}
