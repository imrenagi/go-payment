Subscription API
===

If you want to make recurring payment for your user, you can use Subscription API. For example of usage,
check [Postman Collection](./go-payment.postman_collection.json)

Per current implementation, recurring payment is only supported by Xendit gateway. Here are the recurring payment
features:

1. Supported payment method: Credit Card, Bank Transfer VA, OVO
1. User can make n times recurring payment. If you specify `0` for `total_recurrence`, it creates unlimited number of
   recurring payment.
1. User can `Create`, `Pause`, `Resume` and `Stop` subscription.
1. User can choose its payment method through Xendit Invoice every time they want to pay invoice. If credit card is used
   once, the next payment will be done automatically by charging the credit card.

**Notes:**

1. If user specify `card_token` in subscription create request, the card will directly be charged. You can leave it
   blank in order to use other payment method.

1. If user doesnt pay the latest invoice, the next invoice will still be sent to the user on the next period. If you
   want to automatically stop the subscription when the invoice is not paid, please change this following line
   in [subscription.go](../subscription/subscription.go#L19) to

```
    MissedPaymentAction: MissedPaymentActionStop,
```

> This default value will be configured through configuration file later.

## Create Subscription

### Request

* You can use any integer as `schedule.interval`.
* `day`, `month` and `year` are the valid values for `schedule.interval_unit`
* If `charge_immediately` is `true`, `schedule.start_at` will be ignored and automatically set to now.
* If `charge_immediately` is `false`, `schedule.start_at` must not be empty.

```json
POST /payment/subscriptions

{
	"name": "Imre Subscription",
	"description": "Mantul",
	"amount": 10000,
	"user_id": "imre@gmail.com",
	"currency": "IDR",
	"total_recurrence": 2,
	"card_token": "",
	"charge_immediately": true,
	"schedule": {
		"interval": 1,
		"interval_unit": "day",
		"start_at": "2020-06-10T00:00:00.000Z"
	} 
}
```

### Response

If you are using xendit and charge the user immediately, `last_created_invoice` will not be empty. Use this to redirect
your user to payment page.

```json
{
    "id": 8,    
    "number": "9b9f1999-b2c8-4eb3-8188-23a115237a4c",
    // redacted ....
    // .....
    "charge_immediately": true,
    "last_created_invoice": "https://invoice-staging.xendit.co/web/invoices/5edc5b024eb7c20fe6c8ae91",
    "status": "ACTIVE",
    "recurrence_progress": 0
}
```