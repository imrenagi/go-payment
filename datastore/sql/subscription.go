package sql

import (
	"context"
	"fmt"

	"github.com/imrenagi/go-payment/subscription"

	"github.com/imrenagi/go-payment"
	"gorm.io/gorm"
	"github.com/rs/zerolog"
)

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	r := &SubscriptionRepository{
		DB: db,
	}
	return r
}

type SubscriptionRepository struct {
	DB *gorm.DB
}

func (r SubscriptionRepository) Save(ctx context.Context, subs *subscription.Subscription) error {
	log := zerolog.Ctx(ctx).With().Str("function", "SubscriptionRepository.Save").Logger()

	if err := r.DB.Save(subs).Find(&subs).Error; err != nil {
		log.Error().Err(err).Msg("can't save subscription")
		return payment.ErrDatabase
	}
	return nil
}

func (r *SubscriptionRepository) FindByNumber(ctx context.Context, number string) (*subscription.Subscription, error) {
	log := zerolog.Ctx(ctx).With().
		Str("function", "SubscriptionRepository.FindByNumber").
		Logger()

	var subs subscription.Subscription
	req := r.DB.
		Preload("Schedule").
		Preload("Invoices").
		Where("number = ?", number).First(&subs)

	if req.Error == gorm.ErrRecordNotFound{
		return nil, fmt.Errorf("subscription %s %w", number, payment.ErrNotFound)
	}

	if req.Error != nil {
		log.Error().Err(req.Error).Msg("can't find subscription")
		return nil, payment.ErrDatabase
	}

	return &subs, nil
}
