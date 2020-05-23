package urlvalues

import (
	"fmt"
	"net/url"
	"time"
)

// AddTimeToURLValues append non-zero date time values to URL query values
func AddTimeToURLValues(v *url.Values, t time.Time, fieldName string) {
	if t.IsZero() {
		return
	}
	v.Add(fieldName, t.Format(time.RFC3339))
}

// AddStringSliceToURLValues append values from string slice to URL query values
func AddStringSliceToURLValues(v *url.Values, sl []string, fieldName string) {
	if sl == nil {
		return
	}
	for i, s := range sl {
		v.Add(fmt.Sprintf("%s[%d]", fieldName, i), s)
	}
}
