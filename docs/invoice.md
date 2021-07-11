Invoice API
===

## Generating New Invoice

Use this endpoint to create a payment request to chosen payment channels and gateway.

### Request

```http
POST /payment/invoices
```

```json
{
  "payment": {
    "payment_type": "ovo"
  },
  "customer": {
    "name": "John",
    "email": "foo@example.com",
    "phone_number": "089922222222"
  },
  "items": [{
    "name": "Support Podcast",
    "category": "PODCAST",
    "merchant": "imrenagi.com",
    "description": "donasi podcast imre nagi",
    "qty": 1,
    "price": 80001,
    "currency": "IDR"
  }]
}
```

> To create invoice with `credit_card` payment with/without installment, please take a look [POSTMAN COLLECTION](/example/server/go-payment.postman_collection.json)

### Response

When you call endpoint above, server returns all invoice data. But, to proceed to the payment page you need to pay
attention to `payment` object.

```json
{
  "payment": {
    "id": 48,
    "created_at": "2020-05-25T23:31:44.99873+07:00",
    "updated_at": "2020-05-25T23:31:44.99873+07:00",
    "deleted_at": null,
    "gateway": "xendit",
    "payment_type": "ovo",
    "token": "",
    "redirect_url": "https://invoice.xendit.co/web/invoices5ecbf2f0689543409347ec15",
    "transaction_id": "5ecbf2f0689543409347ec15"
  }
}
```

:heavy_exclamation_mark::heavy_exclamation_mark::heavy_exclamation_mark: Please note:

### For Midtrans Payment Channel

- Value of `payment.gateway` will is always `midtrans`
- You can use `payment.token` to open snap window by using midtrans [snap.js](https://snap-docs.midtrans.com/#snap-js)
- If you want to use [Window Redirection](https://snap-docs.midtrans.com/#window-redirection), you can open a new
  browser tab by using url in `payment.redirect_url`

### For Xendit Payment Channel

- Value of `payment.gateway` will is always `xendit`
- `payment.token` is always empty for all xendit provided payment channels
- You will always open `payment.redirect_url` in new browser tap for all payment methods provided by xendit. Including
  DANA, LinkAja, Kredivo, even Xendit Invoice.
