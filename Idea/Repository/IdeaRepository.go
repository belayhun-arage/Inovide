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

type IdeaRepo struct {
	db *gorm.DB
}

func NewIdeaRepo(sqlite *gorm.DB) Idea.IdeaRepository {
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

func (ideas *IdeaRepo) UpdateIdea(idea *entity.Idea) int64 {
	RowsAffected := ideas.db.Debug().Table("idea").Where("id=? and ideaownerid=?", idea.Id, idea.Ideaownerid).Update(idea).RowsAffected
	// if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
	// 	// handle error
	// }

	defer recover()
	return int64(RowsAffected)

}

func (ideas *IdeaRepo) GetIdea(idea *entity.Idea) int64 {
	numberofAffected := ideas.db.Debug().Table("idea").Where(&entity.Idea{}, idea.Id).Find(idea).RowsAffected
	defer recover()
	return numberofAffected
}

func (ideas *IdeaRepo) DeleteIdea(idea *entity.Idea) int64 {

	RowsAffected := ideas.db.Debug().Table("idea").Where("id=? and ideaownerid =?", idea.Id, idea.Ideaownerid).Delete(idea).RowsAffected
	// if err != nil {
	// 	return err
	// }
	defer recover()
	return RowsAffected
}

func (ideas *IdeaRepo) VoteIdea(ideaid, voterid int) error {

	idea := &entity.Idea{}
	err := ideas.db.Debug().Table("idea").Where("id=?", ideaid).Find(idea).Error

	fmt.Println(idea.Numberofvotes)
	idea.Numberofvotes++
	err = ideas.db.Debug().Table("idea").Where(" id=?", idea.Id).Update(idea).Error

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
func (ideas *IdeaRepo) SearchIdeas(text string, person *entity.Person, searchresults *[]entity.Idea) (*[]entity.Idea, error) {
	var visibility string
	if person.Paid == 0 {
		visibility = "pu"
	}
	err := ideas.db.Table("idea").Debug().Where("visibility=? and title=? ", visibility, text).Find(searchresults).Error

	if err != nil {
		return searchresults, err
	}

	return searchresults, nil
}
