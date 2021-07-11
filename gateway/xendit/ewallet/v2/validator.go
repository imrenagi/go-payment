package ewallet

import (
  "regexp"
)

var (
  ovoChargePhoneRegexp = regexp.MustCompile(`(?:\+62)\d{7,12}`)
)

type ovoPhoneValidator struct {
  pattern *regexp.Regexp
}

func (o ovoPhoneValidator) IsValid(s string) bool {
  return o.pattern.MatchString(s)
}

var OvoChargePhoneValidator = ovoPhoneValidator{pattern: ovoChargePhoneRegexp}
