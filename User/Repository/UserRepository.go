package repository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	"fmt"

	entity "github.com/Samuael/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(sqlite *gorm.DB) *UserRepo {
	return &UserRepo{db: sqlite}
}

func (users *UserRepo) CreateUser(enti *entity.Person) error {

	// dot, err := dotsql.LoadFromFile("sql/userQueries.sql")
	// stri, _ := json.Marshal(enti)

	// fmt.Println(string(stri))

	// stmt, erro := dot.Prepare(users.db, "insert-into-user-table")
	// if erro != nil {
	// 	fmt.Println(erro)
	// 	panic(err)
	// }
	// InsertRessult, err := stmt.Exec(users.db, enti.Firstname, enti.Lastname, enti.Username, enti.Password, enti.Email, enti.Biography, enti.ImageDirectory)

	// if err != nil {
	// 	fmt.Println("Error While inserting the Data  to the table :: User Repository ")
	// 	return err
	// }
	// id, err := InsertRessult.LastInsertId()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// enti.Id = id
	// fmt.Println("Inserted At ::Inside Repository UserRepository ")

	erro := users.db.Table("users").Create(enti).GetErrors()
	fmt.Println("-----------------------")
	if erro != nil {
		fmt.Println("\n\n\nerro\n\n\n\n")
		panic(erro)
	}
	return nil
}

func (users *UserRepo) CheckUser(enti *entity.Person) bool {

	person := entity.Person{}
	users.db.Table("users").Where(&entity.Person{Username: enti.Username, Password: enti.Password}).Find(&person) //Select([]string{"UserName", "Email", "Password"}).Find(person  , )

	fmt.Println(person.Username)
	if person.Username == "" || person.Password == "" || person.Email == "" {
		return false
	}
	return true
}
