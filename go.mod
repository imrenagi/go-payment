module github.com/imrenagi/go-payment

go 1.14

require (
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/rs/cors v1.7.0
	github.com/rs/zerolog v1.18.0
	github.com/spf13/viper v1.6.2
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.6.0
	github.com/veritrans/go-midtrans v0.0.0-20200303064216-54da2d269748
	github.com/xendit/xendit-go v0.8.0
	golang.org/x/sys v0.0.0-20200212091648-12a6c2dcc1e4 // indirect
	gopkg.in/yaml.v2 v2.2.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)

replace github.com/veritrans/go-midtrans v0.0.0-20200303064216-54da2d269748 => github.com/schoters/go-midtrans v0.0.0-20200301123106-412075ea875d
