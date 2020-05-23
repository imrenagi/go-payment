package inmemory

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
)

func NewPaymentConfigRepository() *PaymentConfigRepository {

	var path string = "payment-methods.yml"
	// configFolderPath := os.Getenv("PAYMENT_CONFIG_DIR")
	// if configFolderPath == "" {
	// 	path = "internal/payment/config/payment-methods.yml"
	// } else {
	// 	path = fmt.Sprintf("%s/payment-methods.yml", configFolderPath)
	// }

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error loading creds from path %s : %v", path, err)
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

// PaymentConfigRepository ...
type PaymentConfigRepository struct {
	config *config.PaymentConfig
}

func (r PaymentConfigRepository) FindByPaymentType(
	ctx context.Context,
	paymentType payment.PaymentType,
	opts ...payment.PaymentOption,
) (config.FeeConfigReader, error) {

	options := payment.PaymentOptions{
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
			installment, err := cardPayment.
				GetInstallment(
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
		payment.SourceOtherVA,
		payment.SourceEchannel:
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

func (r PaymentConfigRepository) FindAll(ctx context.Context) (*config.PaymentConfig, error) {
	return r.config, nil
}
