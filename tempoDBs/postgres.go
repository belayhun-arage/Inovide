package postgresSql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	user     = "samuael"
	password = "samuael"
	dbname   = "exercise"
	sslmode  = "disable"
	host     = "localhost"
	port     = 5432
)

var postgresStatmente string
var erro error

func createStatmente() {
	postgresStatmente = fmt.Sprintf("user=%s  password=%s  host=%s  port=%d  dbname=%s sslmode=%s  ", user, password, host, port, dbname, sslmode)

}

func initialize() *sql.DB {
	createStatmente()
	db, erro := sql.Open("postgres", postgresStatmente)

	if erro != nil {
		panic(erro)
	}
	if erro = db.Ping(); erro != nil {
		panic("error While creatinga a database ")
	}

	fmt.Println("Connection Created Succesfully wiith the Database !!!")
	return db
}
