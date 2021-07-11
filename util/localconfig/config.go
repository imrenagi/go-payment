package localconfig

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"fmt"

	"github.com/spf13/viper"
)

// Config used to keep App Config.
type Config struct {
	Test   string `mapstructure:"name"`
	Xendit Xendit `mapstructure:"xendit"`
}

// DBCredential stores database credential
type Xendit struct {
	EWallet EWallet `mapstructure:"ewallet"`
}

type EWallet struct {
	LegacyEnabled bool          `mapstructure:"legacyEnabled"`
	OVO           EWalletConfig `mapstructure:"ovo"`
	Dana          EWalletConfig `mapstructure:"dana"`
	LinkAja       EWalletConfig `mapstructure:"linkaja"`
}

type EWalletConfig struct {
	UseInvoice bool `mapstructure:"invoice"`
	UseLegacy  bool `mapstructure:"legacy"`
}

// LoadConfig reads the file from path and return Secret
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadConfigFromBytes(data)
}

// LoadConfigFromBytes reads the secret file from data bytes
func LoadConfigFromBytes(data []byte) (*Config, error) {
	fang := viper.New()
	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()
	fang.SetEnvPrefix("GOPAYMENT")
	fang.SetConfigType("yaml")

	if err := fang.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	x := fang.Get("name")
	fmt.Println(x)

	var cfg Config
	err := fang.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Error loading creds: %v", err)
	}

	return &cfg, nil
}
