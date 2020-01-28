package Comment

import entity "github.com/Projects/Inovide/models"

type CommentRepo interface {
	CreateComment(comment *entity.Comment) *entity.SystemMessage
	GetComments(comment *[]entity.Comment, id int) int64
	GetCommentsa(id int) []entity.Comment
	UpdateComment(comment *entity.Comment) []error
	DeleteComment(comment *entity.Comment) int64
}
