Payment Methods API
===

## List of payment methods

### Request

If you want to get the estimated admin/installment fee for each payment methods, provice this GET request with
optional `price` and `currency` query. Otherwise, it returns nil `admin_fee` and `installment_fee`

```http
GET /payment/methods?price=1000&currency=IDR
```

### Response

```json
{
  "card_payment": {
    "payment_type": "credit_card",
    "installments": [
      {
        "display_name": "",
        "type": "offline",
        "bank": "bca",
        "terms": [
          {
            "term": 0,
            "admin_fee": {
              "value": 2029,
              "curency": "IDR"
            }
          }
        ]
      }
    ]
  },
  // ... redacted ...
  "ewallets": [
    {
      "payment_type": "gopay",
      "display_name": "Gopay",
      "admin_fee": {
        "value": 0,
        "curency": "IDR"
      }
    }
  ]
}
```
