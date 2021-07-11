package xendit

import (
  "fmt"

  "github.com/xendit/xendit-go/ewallet"
  xinvoice "github.com/xendit/xendit-go/invoice"

  "github.com/imrenagi/go-payment"
  v1 "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v1"
  v2 "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v2"
  "github.com/imrenagi/go-payment/invoice"
)

// NewEWalletChargeRequestFromInvoice create ewallet charge params for xendit ewallet API
func NewEWalletChargeRequestFromInvoice(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {
  switch inv.Payment.PaymentType {
  case payment.SourceOvo:
    return v2.NewOVO(inv)
  case payment.SourceDana:
    return v2.NewDana(inv)
  case payment.SourceLinkAja:
    return v2.NewLinkAja(inv)
  default:
    return nil, fmt.Errorf("unsupported payment method")
  }
}

// Deprecated: NewEwalletRequestFromInvoice creates ewallet request for xendit
func NewEwalletRequestFromInvoice(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {
  switch inv.Payment.PaymentType {
  case payment.SourceOvo:
    return v1.NewOVO(inv)
  case payment.SourceDana:
    return v1.NewDana(inv)
  case payment.SourceLinkAja:
    return v1.NewLinkAja(inv)
  default:
    return nil, fmt.Errorf("payment type is not known")
  }
}

func NewInvoiceRequestFromInvoice(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {

  var reqBuilder invoiceRequestBuilder
  var err error

  req := NewInvoiceRequestBuilder(inv)

  switch inv.Payment.PaymentType {
  case payment.SourceOvo:
    reqBuilder, err = NewOVOInvoice(req)
  case payment.SourceDana:
    reqBuilder, err = NewDanaInvoice(req)
  case payment.SourceLinkAja:
    reqBuilder, err = NewLinkAjaInvoice(req)
  case payment.SourceAlfamart:
    reqBuilder, err = NewAlfamartInvoice(req)
  case payment.SourceBCAVA:
    reqBuilder, err = NewBCAVAInvoice(req)
  case payment.SourceBRIVA:
    reqBuilder, err = NewBRIVAInvoice(req)
  case payment.SourceBNIVA:
    reqBuilder, err = NewBNIVAInvoice(req)
  case payment.SourcePermataVA:
    reqBuilder, err = NewPermataVAInvoice(req)
  case payment.SourceMandiriVA:
    reqBuilder, err = NewMandiriVAInvoice(req)
  case payment.SourceCreditCard:
    reqBuilder, err = NewCreditCardInvoice(req)
  default:
    return nil, fmt.Errorf("payment type is not known")
  }
  if err != nil {
    return nil, err
  }

  return reqBuilder.Build()
}
