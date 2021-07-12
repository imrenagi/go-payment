<!-- markdownlint-disable MD014 MD024 MD026 MD033 MD036 MD041 -->

<div align='center'>

![go-payment](./header.png)

</div>

Payment module used as proxy for multiple payment gateways. Currently it only supports [Midtrans SNAP](https://snap-docs.midtrans.com/) and Xendit [Ewallet](https://xendit.github.io/apireference/#ewallets) and [XenInvoice](https://xendit.github.io/apireference/#invoices). Support for other channels will be added incrementally.

> This payment proxy is a payment service I used for my [personal site](https://imrenagi.com/donate). Thinking that this might be useful to help other people so that they can start accept money ASAP, so I decided to make this module open source.

---

<details>
<summary><b>View table of contents</b></summary><br/>

- [Payment Channels Supported](#payment-channels-supported)
- [Why you should use this payment proxy?](#why-you-should-use-this-payment-proxy)
- [Current Limitations](#current-limitations)
- [Implemented Channels](#implemented-channels)
- [Getting Started](#getting-started)
  - [Payment Gateway Registration](#payment-gateway-registration)
    - [Midtrans](#midtrans)
    - [Xendit](#xendit)
    - [Midtrans VS Xendit Onboarding](#midtrans-vs-xendit-onboarding)
  - [Payment Gateway Callback](#payment-gateway-callback)
    - [Midtrans](#midtrans-1)
    - [Xendit](#xendit-1)
  - [Application Secret](#application-secret)
    - [Database](#database)
    - [Midtrans Credential](#midtrans-credential)
    - [Xendit Credential](#xendit-credential)
  - [Configuration File](#configuration-file)
  - [Mandatory Environment Variables](#mandatory-environment-variables)
- [Example Code](#example-code)
- [API Usage](#api-usage)
- [Contributing](#contributing)
- [License](#license)

</details>

---

## Payment Channels Supported

In general, this payment proxy can support payment through this following channels:

- Recurring payment with Credit Card, OVO, Bank Transfer
- Credit card payment with/without installment
- Ewallet (GoPay, OVO, Dana, LinkAja)
- Retail Outlet (Alfamart, Alfamidi, Dan+Dan)
- Cardless Credit (Akulaku)
- Bank Transfer via Virtual Account (BCA, BNI, BRI, Mandiri, Permata, Other Bank).

> :heavy_exclamation_mark: Recurring payment is only supported via XenditInvoice.

## Why you should use this payment proxy?

- If you are planning to use Midtrans SNAP and Xendit Invoice as the UI for the payment, you are strongly encouraged to use this proxy because it supports both UIs.
- This proxy helps you managing the payment gateway used for each channel. It internally connects to both payment gateway as you need, in no time. What your API user knows is only one single API to generate `Invoice`
- This proxy helps you seemlesly switch the gateway for a payment channel whenever one of them is not functioning properly/down for maintenance. For instance, Bank Transfer by VA, are supported by Midtrans and Xendit. If Midtrans VA is going south, you can easily switch the gateway to Xendit simply by updating the configuration file.
- You can choose whether to absorb the admin/installment fees by yourself or to off load it to your user by changing the payment configuration written in yaml.
- This proxy can generate `Invoice` storing informations about the customer info, item, payment method selected, and its state. `Invoice` state will change over the time depends on the transaction status callback sent by payment gateway.
- You can opt-in to store payment notification callback to your database. Currently it only stores midtrans transaction status. Support for xendit will be added soon.

## Current Limitations
1. For simplify the query creation for database join, I use [gorm.io](https://gorm.io/) as the ORM library. 
1. This proxy is not made for supporting all use cases available out there. It's hard requirement is just so that people can accept payment with as low effort as possible without need to worry about custom UI flow.
1. No callback trigger at least of now once the payment manager is done procesing this request. This will be the next priority of the next release. This issue is documented [here](https://github.com/imrenagi/go-payment/issues/5)
1. Callback or redirect URL is globally configured. This means, you cant configure callback for each request differently on the fly.

## Implemented Channels

This tables shows which payment channels that has been implemented by this proxy.

:white_check_mark: : ready

:heavy_exclamation_mark: : in progress

:x: : not yet supported natively by payment gateway

| Channels                        | Midtrans (Snap)                     | Xendit (ewallet/XenInvoice) |
| ------------------------------- | ----------------------------------- | --------------------------- |
| Credit Card without installment | :white_check_mark:                  | :white_check_mark:          |
| Credit Card with installment    | :white_check_mark:                  | :x:                         |
| BCA VA                          | :white_check_mark:                  | :white_check_mark:          |
| Mandiri VA                      | :white_check_mark:                  | :white_check_mark:          |
| BNI VA                          | :white_check_mark:                  | :white_check_mark:          |
| Permata VA                      | :white_check_mark:                  | :white_check_mark:          |
| Other VA                        | :white_check_mark:                  | :x:                         |
| BRI VA                          | :heavy_exclamation_mark:            | :white_check_mark:          |
| Alfamart, Alfamidi, Dan+Dan     | :white_check_mark:                  | :white_check_mark:          |
| QRIS                            | :white_check_mark:                  | :white_check_mark:          |
| Gopay                           | :white_check_mark:                  | :x:                         |
| OVO                             | :x:                                 | :white_check_mark:          |
| DANA                            | :x:                                 | :white_check_mark:          |
| LinkAja                         | :x:                                 | :white_check_mark:          |
| ShopeePay                       | :white_check_mark:                  | :white_check_mark:          |
| Akulaku                         | :white_check_mark:                  | :x:                         |
| Kredivo                         | :x:                                 | :heavy_exclamation_mark:    |

## Getting Started

Here some preparations that you might need before using this proxy.

### Payment Gateway Registration

This can be tricky. If you have personal business, this might be easier. If you have business entity (PT, CV, etc), there are some additional processes you have to follow and some documents that you have to provide. In this context, I will just assume that you have personal business like what I do: [imrenagi.com](https://imrenagi.com)

#### Midtrans

Please review this [page](https://midtrans.com/tentang-passport) before creating an account.

#### Xendit

Please visit this [registration page](https://dashboard.xendit.co/register/1) for creating an account.

#### Midtrans VS Xendit Onboarding

Here is the comparison between Midtrans and Xendit onboarding based on my onboarding experience.

| Criteria                                      | Midtrans                                                                                             | Xendit                                                                                                 |
| --------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| Document to provide for registration          | KTP, NPWP                                                                                            | KTP, NPWP                                                                                              |
| Cooperation Agreement (Perjanjian Kerja Sama) | Online Signing                                                                                       | Paper Signing and use Legalized Stamp                                                                  |
| Active channels after agreement is signed     | Gopay, Bank Transfer                                                                                 | Bank Transfer, Credit Card                                                                             |
| OVO, LinkAja, Dana activation                 | n/a                                                                                                  | Fill additional forms on the dashboard. Activation varies between weeks - months                       |
| Alfamart activation                           | Require additional request and midtrans review for the activation                                    | No need to be PT, CV. Just fill and sign additional form on the dashboard. Might take weeks or months. |
| Credit card activation                        | Require additional request and midtrans review for the activation                                    | Immediately activated after document sign                                                              |
| Disbursement feature                          | Not included on the same PKS. Need to contact IRIS team for new agreement, activation and onboarding | Immediately activated after document sign                                                              |
| Akulaku activation                            | Might require business entity (PT, CV)                                                               | n/a                                                                                                    |
| Kredivo activation                            | n/a                                                                                                  | Ask your account manager to activate this payment method                                               |
| API Documentation                             | Available                                                                                            | Available                                                                                              |
| Golang SDK                                    | Available                                                                                            | Available, but under development. Expect breaking changes in newer version                             |

### Payment Gateway Callback

Since this library is just providing the http.Handler, you can choose the REST API endpoint used by each callback. 
You can check the example on [server.go](./example/server/server.go).

> Some Xendit Legacy Ewallet API(s) require you to set callback and redirect URL on the body request. You can override this
> value by using environment variables. Go to [Mandatory Environment Variables](#mandatory-environment-variables).

#### Midtrans

To set your callback URL,

- Login to <https://dashboard.midtrans.com>
- Choose environment (Sandbox or Production)
- Click Settings > Configuration
- Set your **Payment Notification URL** with your server callback. For instance: `https://api.imrenagi.com/payment/midtrans/callback`
- Set your **Finish**, **Unfinish**, and **Error** redirect URL
- Click **Update**

#### Xendit

To set your callback URL,

- Login to <https://dashboard.xendit.co>
- Choose environment (Live or Test)
- Click Settings > Callbacks
- Set your callbacks for Invoices Paid. For instance: `https://api.imrenagi.com/payment/xendit/invoice/callback`
- Check option **Also notify my application when an invoice is expired**
- Click **Save and Test**

> LinkAja and DANA callback URL are not defined on xendit dashboard. Instead, they are given while the proxy is initiating the payment request to Xendit API. You can find the callback URL set on [linkaja.go](/gateway/xendit/ewallet/v1/linkaja.go) and [dana.go](/gateway/xendit/ewallet/v1/dana.go)

### Application Secret

Before using this application, you might need to update [secret.yaml](/example/server/secret.yaml) file containing application secret like database and payment gateway credential.

### Application Config

As of now, application config stores configuration about which API that you would like to use for Xendit ewallet payments
such as Dana, OVO, and LinkAja. Please check [config.yaml](/example/server/config.yaml).

```yaml
xendit:
  ewallet:
    ovo:
      invoice: false
      legacy: false
```

* `xendit.ewallet.[ewallet].invoice` set to true if you want to use XenInvoice instead of using direct API integration
* `xendit.ewallet.[ewallet].legacy` set to true if you want to use **legacy** Xendit Ewallet API. Note that this API will be deprecated by first quarter of 2022.

#### Database

I removed MySQL as default database for this library. This library only accept instance of `gorm.DB` for database. 
Thus, you can use any database you like and provide the `gorm.DB` instance of chosen database.

For more, please check [server.go](/example/server/server.go)

#### Midtrans Credential

- Login to <https://dashboard.midtrans.com>
- Choose environment (Sandbox or Production)
- Click Settings > Access Keys
- Grab the credentials, and update the `secret.yaml`

```yaml
payment:
  midtrans:
    secretKey: "midtrans-server-secret"
    clientKey: "midtrans-client-key"
    clientId: "midtrans-merchant-id"
```

#### Xendit Credential

- Login to <https://dashboard.xendit.co>
- Choose environment (Live or Test)
- Click Settings > API Keys > Generate secret key
- Add key Name. Grant write permission for both **Money-in products**
- Take the generated API Keys and Verification Callback Token, update the `secret.yaml`

```yaml
payment:
  ...
  xendit:
    secretKey: "xendit-api-key"
    callbackToken: "xendit-callback-token"
```

### Configuration File

You can take a look sample configuration file named [payment-methods.yml](/example/server/payment-methods.yml). For instance:

```yaml
card_payment:
  payment_type: "credit_card"
  installments:
    - type: offline
      display_name: ""
      gateway: midtrans
      bank: bca
      channel: migs
      default: true
      active: true
      terms:
        - term: 0
          admin_fee:
            IDR:
              val_percentage: 2.9
              val_currency: 2000
              currency: "IDR"
        - term: 3
          installment_fee:
            IDR:
              val_percentage: 5.5
              val_currency: 2200
              currency: "IDR"
```

With above configuration, for installment `offline` with `bca`, you can apply this following fees to the invoice after user generates new invoice:

1. 2.9% + IDR 2000 admin fee for credit card transaction without any installment, or
1. 5.5% + IDR 2200 installment fee for credit card transaction with installment for 3 month tenure.

If you want to absorb the fee, you can simply set `val_percentage` and `val_currency` as `0`

If you only want to apply fee just either by using `val_pecentage` or `val_currency`, simply set the value to one of them and give `0` to the other. For instance:

```yaml
bank_transfers:
  - gateway: midtrans
    payment_type: "bca_va"
    display_name: "BCA"
    admin_fee:
      IDR:
        val_percentage: 0
        val_currency: 4000
        currency: "IDR"
```

> `admin_fee` and `installment_fee` are optional key.

### Mandatory Environment Variables

You need to set these environment variables to make sure this proxy to work.

| Environment Variable  | Required | Description | Example | 
| ------------- | ------------- | ------------- | ------------- |
| ENVIRONMENT  | yes  | decide whether the server is for testing or production. For production, use `prod`.  | `prod`  |
| LOG_LEVEL  | no  | Log level. Default to `DEBUG`. Available values: `DEBUG`, `INFO`, `WARN`, `ERROR`  | `DEBUG`  |
| INVOICE_SUCCESS_REDIRECT_URL | yes | Redirect URL used by xendit if invoice is successfully paid | `http://example.com/donate/thanks` |
| INVOICE_FAILED_REDIRECT_URL | yes | Redirect URL used by xendit if invoice is failed  | `http://example.com/donate/error` |
| DANA_LEGACY_CALLBACK_URL | yes, if you are using legacy ewallet xendit API | Callback URL used for xendit legacy ewallet API to send payment callback | `http://api.example.com/payment/xendit/dana/callback` |
| DANA_LEGACY_REDIRECT_URL | yes, if you are using legacy ewallet xendit API | Redirect URL used by xendit legacy ewallet API to redirect user after payment succeeded  | `http://example.com/donate/thanks` |
| LINKAJA_LEGACY_CALLBACK_URL | yes, if you are using legacy ewallet xendit API | Callback URL used for xendit legacy ewallet API to send payment callback | `http://api.example.com/payment/xendit/linkaja/callback` |
| LINKAJA_LEGACY_REDIRECT_URL | yes, if you are using legacy ewallet xendit API | Redirect URL used by xendit legacy ewallet API to redirect user after payment succeeded  | `http://example.com/donate/thanks` |
| RECURRING_SUCCESS_REDIRECT_URL | yes, if you are using subscription feature | Redirect URL used by xendit subscription API to redirect user after payment succeeded | `http://example.com/donate/thanks` |
| RECURRING_FAILED_REDIRECT_URL | yes, if you are using subscription feature | Redirect URL used by xendit subscription API to redirect user after payment failed | `http://example.com/donate/error` |
| DANA_SUCCESS_REDIRECT_URL | yes, if you are using new xendit ewallet API | Redirect URL used by xendit new ewallet API if payment with dana is success | `http://example.com/success` |
| LINKAJA_SUCCESS_REDIRECT_URL | yes, if you are using new xendit ewallet API | Redirect URL used by xendit new ewallet API if payment with dana is failed | `http://example.com/success` |

## Example Code

You can find the sample code in [here](/example/server)

## API Usage

You can find the details of API usage in [here](/docs)

## Contributing

No rules for now. Feel free to add issue first and optionally submit a PR. Cheers

## License

MIT. Copyright 2022 [Imre Nagi](./LICENSE)
