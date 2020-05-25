package mysql

import (
	"fmt"

	// Adding mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/imrenagi/go-payment/util/localconfig"
	"github.com/jinzhu/gorm"

	// Mysql dialect for gorm
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// URI return mysql uri
func URI(creds localconfig.DBCredential) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
		creds.UserName,
		creds.Password,
		creds.Host,
		creds.Port,
		creds.DBName)
}

//NewGorm return gorm DB connection to mysql
func NewGorm(creds localconfig.DBCredential) *gorm.DB {
	dbURL := URI(creds)
	db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		panic(fmt.Sprintf("gorm cant open connection to mysql: %v", err))
	}

	return db
}
