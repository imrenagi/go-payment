package balance

import (
	"net/url"

	"github.com/xendit/xendit-go"
)

// GetParams contains parameters for Get
type GetParams struct {
	ForUserID   string                        `json:"-"`
	AccountType xendit.BalanceAccountTypeEnum `json:"account_type"`
}

// QueryString creates query string from GetParams, ignores nil values
func (p *GetParams) QueryString() string {
	urlValues := &url.Values{}

	if p.AccountType != "" {
		urlValues.Add("account_type", p.AccountType.String())
	}

	return urlValues.Encode()
}
