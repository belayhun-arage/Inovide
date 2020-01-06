package IdeaRepository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	//entity "github.com/Samuael/Projects/Inovide/models"
	"fmt"

	entity "github.com/Samuael/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type IdeaRepo struct {
	db *gorm.DB
}

func NewIdeaRepo(sqlite *gorm.DB) *IdeaRepo {
	return &IdeaRepo{db: sqlite}
}

func (ideas *IdeaRepo) CreateIdea(idea *entity.Idea) error {
	err := ideas.db.Debug().Table("ideas").Model(&entity.Idea{}).Create(idea).Error
	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		// handle error
	}
	fmt.Println("-----------------------")
	if err != nil {
		fmt.Println("\n\n\n erro \n\n\n\n")
		panic(err)
	}
	return nil
}
