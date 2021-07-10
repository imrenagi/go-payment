package xendit

import (
  "regexp"
)

var (
  // (?:\+62)?0?8\d{2}(\d{5,8})
  // (?:08)\d{2}(\d{5,8})  ==> legacy ovo
  ovoPhoneRegexp = regexp.MustCompile(`(?:08)\d{7,12}`)
  ovoChargePhoneRegexp = regexp.MustCompile(`(?:\+62)\d{7,12}`)
)

type ovoPhoneValidator struct {
  pattern *regexp.Regexp
}

func (o ovoPhoneValidator) IsValid(s string) bool {
  return o.pattern.MatchString(s)
}

var OvoChargePhoneValidator = ovoPhoneValidator{pattern: ovoChargePhoneRegexp}
var OvoPhoneValidator = ovoPhoneValidator{pattern: ovoPhoneRegexp}
