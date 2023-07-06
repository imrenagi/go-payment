package localconfig

import (
	"bytes"

	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Secret used to keep App Secret.
type Secret struct {
	Setting SettingSecret `yaml:"setting"`
	DB      DBCredential  `yaml:"db"`
	Payment PaymentSecret `yaml:"payment"`
}

// SettingSecret stores config setting on secret yaml
type SettingSecret struct {
	RunningPort string `yaml:"runningPort"`
}

// DBCredential stores database credential
type DBCredential struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Timezone string `yaml:"timezone"`
}

// PaymentSecret stores secret for payment gateway
type PaymentSecret struct {
	Midtrans APICredential `yaml:"midtrans"`
	Xendit   APICredential `yaml:"xendit"`
}

// APICredential stores the credential used for connecting to an API service
type APICredential struct {
	ClientID      string `yaml:"clientId"`
	ClientKey     string `yaml:"clientKey"`
	SecretKey     string `yaml:"secretKey"`
	CallbackToken string `yaml:"callbackToken"`
}

// LoadSecret reads the file from path and return Secret
func LoadSecret(path string) (*Secret, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadSecretFromBytes(data)
}

// LoadSecretFromBytes reads the secret file from data bytes
func LoadSecretFromBytes(data []byte) (*Secret, error) {
	fang := viper.New()
	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()
	fang.SetEnvPrefix("GOPAYMENT")
	fang.SetConfigType("yaml")

	if err := fang.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	var creds Secret
	err := fang.Unmarshal(&creds)
	if err != nil {
		log.Fatalf("Error loading creds: %v", err)
	}

	return &creds, nil
}
