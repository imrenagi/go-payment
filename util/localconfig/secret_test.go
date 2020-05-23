package localconfig_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/imrenagi/go-payment/util/localconfig"
	"github.com/stretchr/testify/assert"
)

var cfg = []byte(`
db:
  host: "127.0.0.1"
  port: 3306
  username: "john"
  password: "passwd"
  dbname: "dbName"
payment:
  midtrans:
    secretKey: "midtranssecretkey"
    clientId: "midtransclientid"
`)

func TestLoadSecret(t *testing.T) {
	err := ioutil.WriteFile("/tmp/secret.yaml", cfg, 0644)
	if err != nil {
		t.Fail()
		t.Logf("Can't prepare temporary file for test. Got `%v`", err)
	}

	t.Run("Test read config from a file", func(t *testing.T) {
		secret, _ := LoadSecret("/tmp/secret.yaml")
		assert.Equal(t, "127.0.0.1", secret.DB.Host)
		assert.Equal(t, 3306, secret.DB.Port)
		assert.Equal(t, "john", secret.DB.UserName)
		assert.Equal(t, "passwd", secret.DB.Password)
		assert.Equal(t, "dbName", secret.DB.DBName)
		assert.Equal(t, "midtransclientid", secret.Payment.Midtrans.ClientID)
		assert.Equal(t, "midtranssecretkey", secret.Payment.Midtrans.SecretKey)
	})

	t.Run("Test env var overriding", func(t *testing.T) {
		os.Setenv("GOPAYMENT_DB_HOST", "localhost")
		secret, _ := LoadSecret("/tmp/secret.yaml")
		assert.Equal(t, "localhost", secret.DB.Host)
		assert.Equal(t, 3306, secret.DB.Port)
		assert.Equal(t, "john", secret.DB.UserName)
		assert.Equal(t, "passwd", secret.DB.Password)
		assert.Equal(t, "dbName", secret.DB.DBName)
	})

	err = os.Remove("/tmp/secret.yaml")
	if err != nil {
		t.Fail()
		t.Logf("Can't clean up test resource")
	}
}
