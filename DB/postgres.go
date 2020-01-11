package config

import (
	"fmt"

	// _ "github.com/lib/pq"
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var dbs *sql.DB

var postgresStatmente string
var errors error

const (
	username = "postgres"
	password = "faniman093864"
	host     = "localhost"
	dbname   = "inovide"
)

func InitializPosts() (*sql.DB, error) {
	postgresStatmente = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, dbname)

	dbs, errors = sql.Open("postgres", postgresStatmente)

	if errors != nil {
		panic(errors)

	}

	if errors = dbs.Ping(); errors != nil {
		panic(errors)

	}
	fmt.Println("Succesfull Registration ")

	return dbs, nil

}

func InitializPostgres() (*gorm.DB, error) {
	// Preparing the statmente
	postgresStatmente = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, dbname)

	db, errors = gorm.Open("postgres", postgresStatmente)

	if errors != nil {
		panic(errors)

	}

	fmt.Println("Succesfull Registration ")

	return db, nil

}
