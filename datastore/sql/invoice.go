package sql

import (
	"context"
	"fmt"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	"gorm.io/gorm"
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
		Preload("LineItems").
		Preload("BillingAddress").
		Where("number = ?", number).Find(&invoice)

	if req.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("invoice %s %w", number, payment.ErrNotFound)
	}

	if req.Error != nil {
		log.Error().Err(req.Error).Msg("can't find invoice")
		return nil, payment.ErrDatabase
	}
	return &invoice, nil
}
