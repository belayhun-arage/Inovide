package IdeaRepository

/*This package Will Responsibel For Manipulating the database and Genereating an Instance of User to be used by the Service */
import (
	//entity "github.com/Samuael/Projects/Inovide/models"
	"fmt"

	entity "github.com/Projects/Inovide/models"
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

func (ideas *IdeaRepo) GetIdea(id int) (*entity.Idea, error) {

	idea := &entity.Idea{}
	err := ideas.db.Debug().Table("idea").Find(idea, id).Error

	if err != nil {
		return nil, err
	}
	return idea, nil

}
func (ideas *IdeaRepo) DeleteIdea(id int) error {

	err := ideas.db.Debug().Table("idea").Delete(&entity.Person{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (ideas *IdeaRepo) VoteIdea(ideaid, voterid int) error {

	idea := &entity.Idea{}
	err := ideas.db.Debug().Table("idea").Where(entity.Idea{}, ideaid).Find(idea).Error

	fmt.Println(idea.Numberofvotes)
	err = ideas.db.Debug().Table("idea").Where(" id=?", idea.Id).Update(&entity.Idea{Ideaownerid: idea.Ideaownerid, Title: idea.Title, Description: idea.Description, Visibility: idea.Visibility, Numberofvotes: idea.Numberofvotes + 1, Numberofcomment: idea.Numberofcomment}).Error

	if err != nil {
		return err
	}

	votee := &entity.Votee{}
	votee.Ideaid = ideaid
	votee.Voterid = voterid
	err = ideas.db.Table("votee").Save(votee).Error
	if err != nil {
		return err
	}
	return nil
}

// func (ideas *IdeaRepo) CheckVoteIdea(ideaid, voterid int) error {

// 	return nil
// }
