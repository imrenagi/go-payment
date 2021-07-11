package ewallet

import (
  "regexp"
)

var (
  ovoPhoneRegexp = regexp.MustCompile(`(?:08)\d{7,12}`)
)

type ovoPhoneValidator struct {
  pattern *regexp.Regexp
}

func (o ovoPhoneValidator) IsValid(s string) bool {
  return o.pattern.MatchString(s)
}

var OvoPhoneValidator = ovoPhoneValidator{pattern: ovoPhoneRegexp}
