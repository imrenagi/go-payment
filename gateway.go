package payment

import (
	"fmt"
	"strings"
)

type Gateway int

const (
	UnknownGateway Gateway = iota
	Midtrans
	Xendit
)

func (g Gateway) String() string {
	return []string{"unkown", "midtrans", "xendit"}[g]
}

func (g *Gateway) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var n string
	if err := unmarshal(&n); err != nil {
		return err
	}

	*g = NewGateway(n)

	if *g == UnknownGateway {
		return fmt.Errorf("payment gateway is not recognized")
	}

	return nil
}

func NewGateway(name string) Gateway {
	var g Gateway
	switch strings.ToLower(name) {
	case "midtrans":
		g = Midtrans
	case "xendit":
		g = Xendit
	default:
		g = UnknownGateway
	}
	return g
}
