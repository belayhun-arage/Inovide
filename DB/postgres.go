package config

import (
	"fmt"

	// _ "github.com/lib/pq"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var postgresStatmente string
var errors error

const (
	username = "samuael"
	password = "samuaelfirst"
	host     = "localhost"
	dbname   = "inovide"
)

func InitializPostgres() (*gorm.DB, error) {
	// Preparing the statmente
	postgresStatmente = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, dbname)

	db, errors = gorm.Open("postgres", postgresStatmente)

	if errors != nil {
		panic(errors)

	}

	// if errors = db.Ping(); errors != nil {
	// 	panic(errors)
	// 	//  if i write a code below the panic statmente it will be Unreachable

	// }
	fmt.Println("Succesfull Registration ")

	return db, nil

}
