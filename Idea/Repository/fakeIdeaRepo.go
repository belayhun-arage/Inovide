package IdeaRepository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	//entity "github.com/Samuael/Projects/Inovide/models"
	"fmt"

	"github.com/Projects/Inovide/Idea"
	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type FakeIdeaRepo struct {
	db *gorm.DB
}

func NewFakeIdeaRepo(sqlite *gorm.DB) Idea.FakeIdeaRepository {
	return &FakeIdeaRepo{db: sqlite}
}

func (ideas *FakeIdeaRepo) CreateFakeIdea(idea *entity.Idea) error {
	err := ideas.db.Debug().Table("idea").Model(&entity.Idea{}).Create(idea).Error
	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		// handle error
	}
	fmt.Println("-----------------------")
	if err != nil {
		fmt.Println("   erro ")
		panic(err)
	}
	return nil
}
