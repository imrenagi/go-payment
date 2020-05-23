package midtrans

import "strings"

// EnvironmentType value
type EnvironmentType int8

const (
	_ EnvironmentType = iota

	// Sandbox : represent sandbox environment
	Sandbox

	// Production : represent production environment
	Production
)

var typeString = map[EnvironmentType]string{
	Sandbox:    "https://api.sandbox.midtrans.com",
	Production: "https://api.midtrans.com",
}

// implement stringer
func (e EnvironmentType) String() string {
	for k, v := range typeString {
		if k == e {
			return v
		}
	}
	return "undefined"
}

// SnapURL : Get environment API URL
func (e EnvironmentType) SnapURL() string {
	return strings.Replace(e.String(), "api.", "app.", 1)
}

// IrisURL : Get environment API URL
func (e EnvironmentType) IrisURL() string {
	return strings.Replace(e.String(), "api.", "app.", 1) + "/iris"
}
