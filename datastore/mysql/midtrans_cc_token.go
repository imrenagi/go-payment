package mysql

import (
	"context"
	"fmt"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

type MidtransCardTokenRepository struct {
	DB *gorm.DB
}

func (r MidtransCardTokenRepository) Save(ctx context.Context, token *midtrans.CardToken) error {

	log := zerolog.Ctx(ctx).With().Str("function", "MidtransCardTokenRepository.Save").Logger()

	if err := r.DB.Save(token).Find(&token).Error; err != nil {
		log.Error().Err(err).Msgf("cant save midtrans card token")
		return payment.ErrDatabase
	}
	return nil
}

func (r MidtransCardTokenRepository) FindAllByUserID(ctx context.Context, userID string) ([]midtrans.CardToken, error) {
	log := zerolog.Ctx(ctx).With().Str("function", "MidtransCardTokenRepository.FindAllByUserID").Logger()

	var tokens []midtrans.CardToken
	req := r.DB.
		Where("user_id = ?", userID).
		Find(&tokens)

	if req.RecordNotFound() {
		return nil, fmt.Errorf("card token for userID %s %w", userID, payment.ErrNotFound)
	}
	errs := req.GetErrors()
	if len(errs) > 0 {
		log.Error().Err(errs[0]).Msg("cant find midtrans card token")
		return nil, payment.ErrDatabase
	}
	return tokens, nil
}
