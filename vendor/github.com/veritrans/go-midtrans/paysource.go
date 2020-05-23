package midtrans

// PaymentType value
type PaymentType string

const (
	// SourceBankTransfer : bank_transfer
	SourceBankTransfer PaymentType = "bank_transfer"

	// SourceBNIVA : bni_va
	SourceBNIVA PaymentType = "bni_va"

	// SourcePermataVA : permata_va
	SourcePermataVA PaymentType = "permata_va"

	// SourceBCAVA : bca_va
	SourceBCAVA PaymentType = "bca_va"

	// SourceOtherVA : other_va
	SourceOtherVA PaymentType = "other_va"

	// SourceBcaKlikpay : bca_klikpay
	SourceBcaKlikpay PaymentType = "bca_klikpay"

	// SourceBriEpay : bri_epay
	SourceBriEpay PaymentType = "bri_epay"

	// SourceCreditCard : credit_card
	SourceCreditCard PaymentType = "credit_card"

	// SourceCimbClicks : cimb_clicks
	SourceCimbClicks PaymentType = "cimb_clicks"

	// SourceDanamonOnline : danamon_online
	SourceDanamonOnline PaymentType = "danamon_online"

	// SourceConvStore : cstore
	SourceConvStore PaymentType = "cstore"

	// SourceKlikBca : bca_klikbca
	SourceKlikBca PaymentType = "bca_klikbca"

	// SourceEchannel : echannel
	SourceEchannel PaymentType = "echannel"

	// SourceMandiriClickpay : mandiri_clickpay
	SourceMandiriClickpay PaymentType = "mandiri_clickpay"

	// SourceTelkomselCash : telkomsel_cash
	SourceTelkomselCash PaymentType = "telkomsel_cash"

	// SourceIndosatDompetku : indosat_dompetku
	SourceIndosatDompetku PaymentType = "indosat_dompetku"

	// SourceMandiriEcash : mandiri_ecash
	SourceMandiriEcash PaymentType = "mandiri_ecash"

	// SourceKioson : kioson
	SourceKioson PaymentType = "kioson"

	// SourceIndomaret : indomaret
	SourceIndomaret PaymentType = "indomaret"

	// SourceAlfamart : alfamart
	SourceAlfamart PaymentType = "alfamart"

	// SourceGiftCardIndo : gci
	SourceGiftCardIndo PaymentType = "gci"

	// SourceGopay : gopay
	SourceGopay PaymentType = "gopay"

	// SourceAkulaku : akulaku
	SourceAkulaku PaymentType = "akulaku"
)

// AllPaymentSource : Get All available PaymentType
var AllPaymentSource = []PaymentType{
	SourceGopay,
	SourceCreditCard,
	SourceMandiriClickpay,
	SourceCimbClicks,
	SourceDanamonOnline,
	SourceKlikBca,
	SourceBcaKlikpay,
	SourceBriEpay,
	SourceTelkomselCash,
	SourceEchannel,
	SourceIndosatDompetku,
	SourceMandiriEcash,
	SourceBNIVA,
	SourcePermataVA,
	SourceBCAVA,
	SourceIndomaret,
	SourceKioson,
	SourceGiftCardIndo,
}
