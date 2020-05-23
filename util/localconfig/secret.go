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
	DB      DBCredential  `yaml:"db"`
	Payment PaymentSecret `yaml:"payment"`
}

type DBCredential struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type PaymentSecret struct {
	Midtrans APICredential `yaml:"midtrans"`
	Xendit   APICredential `yaml:"xendit"`
}

type APICredential struct {
	ClientID      string `yaml:"clientId"`
	ClientKey     string `yaml:"clientKey"`
	SecretKey     string `yaml:"secretKey"`
	CallbackToken string `yaml:"callbackToken"`
}

func LoadSecret(path string) (*Secret, error) {
	fang := viper.New()
	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()
	fang.SetEnvPrefix("GOPAYMENT")
	fang.SetConfigType("yaml")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := fang.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	var creds Secret
	err = fang.Unmarshal(&creds)
	if err != nil {
		log.Fatalf("Error loading creds from path %s : %v", path, err)
	}

	return &creds, nil
}
