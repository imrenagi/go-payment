package manager

import (
	"context"
	"errors"
	"time"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/rs/zerolog"
	midgo "github.com/veritrans/go-midtrans"
)

// ProcessMidtransCallback takes care of notification sent by midtrans. This checks the validity of the sign key and the similarity
// between the notification and transaction satus.
func (m *Manager) ProcessMidtransCallback(ctx context.Context, mr midgo.Response) error {

	log := zerolog.Ctx(ctx).
		With().
		Str("function", "PaymentService.ProcessMidtransCallback()").
		Str("cmd_order_id", mr.OrderID).
		Str("cmd_transaction_id", mr.TransactionID).
		Str("cmd_gross_amount", mr.GrossAmount).
		Str("cmd_transaction_status", mr.TransactionStatus).
		Logger()

	var storedStatus *midtrans.TransactionStatus
	var err error

	if m.midTransactionRepository != nil {
		storedStatus, err = m.midTransactionRepository.FindByOrderID(ctx, mr.OrderID)
		if err != nil && !errors.Is(err, payment.ErrNotFound) {
			return err
		}
	}

	ttt, err := time.Parse("2006-01-02 15:04:05", mr.TransactionTime)
	if err != nil {
		log.Error().Err(err).Msg("cant parse transaction time")
		return payment.ErrInternal
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	transactionTime := ttt.In(loc)

	if storedStatus == nil {
		storedStatus = &midtrans.TransactionStatus{
			StatusCode:        mr.StatusCode,
			StatusMessage:     mr.StatusMessage,
			SignKey:           mr.SignKey,
			Bank:              mr.Bank,
			FraudStatus:       mr.FraudStatus,
			PaymentType:       mr.PaymentType,
			OrderID:           mr.OrderID,
			TransactionID:     mr.TransactionID,
			TransactionTime:   transactionTime,
			TransactionStatus: mr.TransactionStatus,
			GrossAmount:       mr.GrossAmount,
			MaskedCard:        mr.MaskedCard,
			Currency:          mr.Currency,
			CardType:          mr.CardType,
		}

	} else {
		storedStatus.StatusCode = mr.StatusCode
		storedStatus.StatusMessage = mr.StatusMessage
		storedStatus.GrossAmount = mr.GrossAmount
		storedStatus.FraudStatus = mr.FraudStatus
		storedStatus.SignKey = mr.SignKey
		storedStatus.TransactionTime = transactionTime
		storedStatus.TransactionStatus = mr.TransactionStatus
		storedStatus.TransactionID = mr.TransactionID
		storedStatus.PaymentType = mr.PaymentType
		storedStatus.MaskedCard = mr.MaskedCard
		storedStatus.CardType = mr.CardType
		storedStatus.Bank = mr.Bank
	}

	if err := storedStatus.IsValid(m.midtransGateway.NotificationValidationKey()); err != nil {
		return err
	}

	if m.midTransactionRepository != nil {
		err = m.midTransactionRepository.Save(ctx, storedStatus)
		if err != nil {
			return err
		}
	}

	err = m.processNotification(ctx, *storedStatus)
	if err != nil {
		log.Error().Err(err).Msg("something wrong when publishing")
		return err
	}

	return nil
}

func (m *Manager) processNotification(ctx context.Context, status midtrans.TransactionStatus) error {

	log := zerolog.Ctx(ctx).With().
		Str("transaction_status", status.TransactionStatus).
		Str("payment_type", status.PaymentType).
		Str("fraud_status", status.FraudStatus).
		Logger()

	switch status.TransactionStatus {
	case "capture":
		if status.PaymentType == "credit_card" && status.FraudStatus == "accept" {

			_, err := m.PayInvoice(ctx, &PayInvoiceRequest{
				InvoiceNumber: status.OrderID,
				TransactionID: status.TransactionID,
			})
			if err != nil {
				return err
			}
		} else {
			log.Warn().Msg("transaction captured, potentially fraud")
			return nil
		}
	case "settlement":
		_, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: status.OrderID,
			TransactionID: status.TransactionID,
		})
		if err != nil {
			return err
		}
	case "deny", "expire", "cancel":
		_, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			TransactionID: status.TransactionID,
			InvoiceNumber: status.OrderID,
			Reason:        status.TransactionStatus,
		})
		if err != nil {
			return err
		}
	case "pending":
		_, err := m.ProcessInvoice(ctx, status.OrderID)
		if err != nil {
			return err
		}
	default:
		log.Warn().Msg("payment status type is unidentified")
		return nil
	}

	return nil
}
