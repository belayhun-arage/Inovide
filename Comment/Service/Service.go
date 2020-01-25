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

func (commentservice *CommentService) UpdateComment(comment *entity.Comment) *entity.SystemMessage {

	systemmessage := &entity.SystemMessage{}

	err := commentservice.CommentRepo.UpdateComment(comment)
	if err != nil {
		systemmessage.Succesful = false
		systemmessage.Message = "Can't Update the Comment "
	} else {
		systemmessage.Succesful = true
		systemmessage.Message = "The Comment is Updated "
	}
	return systemmessage
}
func (commentservice *CommentService) DeleteComment(comment *entity.Comment) *entity.SystemMessage {

	systetemmessage := &entity.SystemMessage{}
	erro := commentservice.CommentRepo.DeleteComment(comment)
	if erro != nil {

		systetemmessage.Message = "Can't Delete the Comment System Error "
		systetemmessage.Succesful = false
	} else {
		systetemmessage.Message = "Succesfully Deleted The Comment "
		systetemmessage.Succesful = true
	}
	return systetemmessage
}
