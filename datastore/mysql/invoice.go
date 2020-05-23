package mysql

import (
	"context"
	"fmt"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	r := &InvoiceRepository{
		DB: db,
	}

	return r
}

type InvoiceRepository struct {
	DB *gorm.DB
}

func (r InvoiceRepository) Save(ctx context.Context, invoice *invoice.Invoice) error {
	log := zerolog.Ctx(ctx).With().Str("function", "InvoiceRepository.Save").Logger()

	if err := r.DB.Save(invoice).Find(&invoice).Error; err != nil {
		log.Error().Err(err).Msg("can't save invoice")
		return payment.ErrDatabase
	}
	return nil
}

func (r *InvoiceRepository) FindByNumber(ctx context.Context, number string) (*invoice.Invoice, error) {
	log := zerolog.Ctx(ctx).With().
		Str("function", "InvoiceRepository.FindByNumber").
		Logger()

	var invoice invoice.Invoice
	req := r.DB.
		Preload("Payment").
		Preload("Payment.CreditCardDetail").
		Preload("LineItem").
		Preload("BillingAddress").
		Where("number = ?", number).Find(&invoice)

	if req.RecordNotFound() {
		return nil, fmt.Errorf("invoice %s %w", number, payment.ErrNotFound)
	}

	errs := req.GetErrors()
	if len(errs) > 0 {
		log.Error().Err(errs[0]).Msg("can't find invoice")
		return nil, payment.ErrDatabase
	}
	return &invoice, nil
}
