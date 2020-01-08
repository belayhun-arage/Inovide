package repository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	"fmt"

	entity "github.com/Samuael/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(sqlite *gorm.DB) *UserRepo {
	return &UserRepo{db: sqlite}
}

func (users *UserRepo) CreateUser(enti *entity.Person) error {

	err := users.db.Debug().Table("users").Model(&entity.Person{}).Create(enti).Error
	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		// handle error
	}
	fmt.Println("-----------------------")
	if err != nil {
		return err
	}
	return nil
}

func (users *UserRepo) CheckUser(enti *entity.Person) bool {

	person := entity.Person{}
	users.db.Table("users").Select("ID").Debug().Model(&entity.Person{}).Where("UserName=$1 AND Password=$2", enti.Username, enti.Password).Find(&person) //Select([]string{"UserName", "Email", "Password"}).Find(person  , )

	fmt.Println(person.ID, "_______-------<< User Repo")
	// fmt.Println(peoples.Username, peoples.Password, peoples.Email)
	if person.Username == "" || person.Password == "" || person.Email == "" {
		return false
	}
	return true
}
