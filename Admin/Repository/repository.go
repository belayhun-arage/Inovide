package AdminRepository

import (
	"fmt"

	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

type AdminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(dbs *gorm.DB) *AdminRepo {
	return &AdminRepo{db: dbs}
}

func (adminrepo *AdminRepo) CountUsers() int64 {

	count := adminrepo.db.Table("users").Find([]entity.Person{}).RowsAffected
	fmt.Println(count)

	return count

}

func (adminrepo *AdminRepo) CountIdeas() int64 {

	count := adminrepo.db.Table("idea").Find([]entity.Idea{}).RowsAffected
	fmt.Println(count)
	return count

}
func (adminrepo *AdminRepo) CountAdmins() int64 {

	// val := true
	count := adminrepo.db.Table("users").Find([]entity.Person{}).RowsAffected
	fmt.Println(count)

	return count

}
func (adminrepo *AdminRepo) CountMessages() int64 {

	count := adminrepo.db.Table("message").Find([]entity.Message{}).RowsAffected
	fmt.Println(count)

	return count

}
func (adminrepo *AdminRepo) CountActiveUsers() int64 {

	count := adminrepo.db.Table("session").Find([]entity.Session{}).RowsAffected
	fmt.Println(count)

	return count

}
