package CommentRepo

import (
	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(dbs *gorm.DB) *CommentRepo {
	return &CommentRepo{db: dbs}
}

func (commentrepo *CommentRepo) CreateComment(comment *entity.Comment) error {
	err := commentrepo.db.Table("comment").Debug().Model(&entity.Comment{}).Save(comment).Error
	if err != nil {
		return err
	}
	return nil
}

func (commentrepo *CommentRepo) GetComments(comment *[]entity.Comment, id int) error {
	err := commentrepo.db.Table("comment").Where(&entity.Comment{}, id).Debug().Find(comment).Error
	if err != nil {
		return err
	}
	return nil
}

func (commentrepo *CommentRepo) DeleteComment(commentid int) error {

	return nil
}
