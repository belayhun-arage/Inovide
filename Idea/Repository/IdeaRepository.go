package repository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	//entity "github.com/Samuael/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

type IdeaRepo struct {
	db *gorm.DB
}

func NewIdeaRepo(sqlite *gorm.DB) *IdeaRepo {
	return &IdeaRepo{db: sqlite}
}
