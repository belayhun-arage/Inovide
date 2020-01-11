package CommentService

import (
	CommentRepo "github.com/Projects/Inovide/Comment/Repository"
	entity "github.com/Projects/Inovide/models"
)

type CommentService struct {
	CommentRepo *CommentRepo.CommentRepo
}

func NewCommentService(commentrepo *CommentRepo.CommentRepo) *CommentService {
	return &CommentService{CommentRepo: commentrepo}
}

func (commentservice *CommentService) CreateComment(comment *entity.Comment) *entity.SystemMessage {
	systemMessage := &entity.SystemMessage{}

	err := commentservice.CommentRepo.CreateComment(comment)
	if err != nil {
		systemMessage.Succesful = false
		systemMessage.Message = "Can't Save The Comment "
		return systemMessage
	}
	systemMessage.Message = "SuCcesfully inserted the message "
	systemMessage.Succesful = true
	return systemMessage
}

func (commentservice *CommentService) GetComments(comment *[]entity.Comment, id int) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}

	err := commentservice.CommentRepo.GetComments(comment, id)
	if err != nil {
		systemmessage.Message = "Can't et The Messages in the specified Id "
		systemmessage.Succesful = true
	}
	systemmessage.Message = "Succesfully Fetched The Comments "
	systemmessage.Succesful = true
	return systemmessage

}
