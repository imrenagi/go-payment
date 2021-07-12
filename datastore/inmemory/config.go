package inmemory

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
)

// NewPaymentConfigRepository creates payment configuration data source by reading
// a file located on `source`
func NewPaymentConfigRepository(source string) *PaymentConfigRepository {

	data, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalf("Error loading creds from path %s : %v", source, err)
		return nil
	}

	config, err := config.LoadPaymentConfigs(data)
	if err != nil {
		return nil
	}

	repo := &PaymentConfigRepository{
		config: config,
	}

	return repo
}

// PaymentConfigRepository is storage for payment configurations
type PaymentConfigRepository struct {
	config *config.PaymentConfig
}

// FindByPaymentType return FeeConfigReader for a given payment type. If it is a credit card,
// credit card option will be check to get the type of installment, term and its aqcuiring bank.
// Otherwise 0 month installment offline from BCA will be used.
func (r PaymentConfigRepository) FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.Option) (config.FeeConfigReader, error) {

	options := payment.Options{
		CreditCard: &payment.CreditCard{
			Bank: payment.BankBCA,
			Installment: payment.Installment{
				Type: payment.InstallmentOffline,
				Term: 0,
			},
		},
	}

	for _, o := range opts {
		o(&options)
	}

	switch paymentType {
	case payment.SourceCreditCard:
		cardPayment := r.config.CardPayment
		if options.CreditCard != nil {
			installment, err := cardPayment.GetInstallment(
				options.CreditCard.Bank,
				options.CreditCard.Installment.Type,
			)
			if err != nil {
				return nil, err
			}

			term, err := installment.GetTerm(options.CreditCard.Installment.Term)
			if err != nil {
				return nil, err
			}

			return term, nil
		}
	case payment.SourceBCAVA,
		payment.SourcePermataVA,
		payment.SourceBNIVA,
		payment.SourceBRIVA,
		payment.SourceOtherVA,
		payment.SourceMandiriVA:
		for _, bt := range r.config.BankTransfers {
			if bt.PaymentType == paymentType {
				return &bt, nil
			}
		}
	case payment.SourceAlfamart:
		for _, cstore := range r.config.CStores {
			if cstore.PaymentType == paymentType {
				return &cstore, nil
			}
		}
	case payment.SourceGopay,
		payment.SourceOvo,
		payment.SourceLinkAja,
		payment.SourceShopeePay,
		payment.SourceQRIS,
		payment.SourceDana:
		for _, ewallet := range r.config.EWallets {
			if ewallet.PaymentType == paymentType {
				return &ewallet, nil
			}
		}
	case payment.SourceAkulaku:
		for _, card := range r.config.CardlessCredits {
			if card.PaymentType == paymentType {
				return &card, nil
			}
		}
	default:
		return nil, fmt.Errorf("payment type %w", payment.ErrNotFound)
	}
	return nil, nil
}

// FindAll returns all payment configurations.
func (r PaymentConfigRepository) FindAll(ctx context.Context) (*config.PaymentConfig, error) {
	// TODO we need to check whether the payment method is enabled or not
	return r.config, nil
}
