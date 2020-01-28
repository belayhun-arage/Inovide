package Comment

import (
	entity "github.com/Projects/Inovide/models"
)

type CommentService interface {
	CreateComment(comment *entity.Comment) *entity.SystemMessage
	GetComments(comment *[]entity.Comment, id int) *entity.SystemMessage
	UpdateComment(comment *entity.Comment) *entity.SystemMessage
	DeleteComment(comment *entity.Comment) *entity.SystemMessage
}
