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

func (commentrepo *CommentRepo) GetComments(comment *[]entity.Comment, id int) int64 {
	err := commentrepo.db.Table("comment").Where("ideaid=?", id).Debug().Find(comment).RowsAffected
	defer recover()
	return err
}

func (commentrepo *CommentRepo) GetCommentsa(id int) []entity.Comment {
	comments := []entity.Comment{}
	bd := commentrepo.db.Table("comment").Where("ideaid=?", id).Find(&comments).RowsAffected
	defer recover()
	if bd >= 1 {
		return comments
	} else {
		return []entity.Comment{}
	}
}
func (CommentRepo *CommentRepo) UpdateComment(comment *entity.Comment) []error {
	erro := CommentRepo.db.Model(&entity.Comment{}).Table("comment").Save(comment).GetErrors()
	return erro
}
func (commentrepo *CommentRepo) DeleteComment(comment *entity.Comment) int64 {
	affecteds := commentrepo.db.Table("comment").Debug().Where("id=? and commentorid=?", comment.Id,
		comment.Commentorid).Delete(comment).RowsAffected
	defer recover()
	return affecteds
}
