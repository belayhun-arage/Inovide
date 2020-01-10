package ideaService

import (
	IdeaRepository "github.com/Projects/Inovide/Idea/Repository"
	entity "github.com/Projects/Inovide/models"
)

type IdeaService struct {
	Idearepo *IdeaRepository.IdeaRepo
}

func NewIdeaService(idearep *IdeaRepository.IdeaRepo) *IdeaService {
	return &IdeaService{Idearepo: idearep}
}
func (ideaServise *IdeaService) CreateIdea(idea *entity.Idea) *entity.SystemMessage {
	var message = entity.SystemMessage{}
	if idea.Title == "" {
		message.Message = "The Input Is Not Fully FIlled Please Submitt Again Filling The Data Appropriately"
		message.Succesful = false
	} else {
		er := ideaServise.Idearepo.CreateIdea(idea)
		if er != nil {
			panic(er)
		}
	}
	message.Message = "Succesfully Inserted "
	message.Succesful = true
	return &message
}

func (ideaServise *IdeaService) GetIdea(idea *entity.Idea, id int) (*entity.Idea, *entity.SystemMessage) {
	systemMessage := &entity.SystemMessage{}

	ideaFromRepo, erro := ideaServise.Idearepo.GetIdea(id)

	if erro != nil {
		systemMessage.Message = "No Can't Get  the Idea "
		systemMessage.Succesful = false
		return nil, systemMessage
	}
	systemMessage.Message = "Succesfully Found "
	systemMessage.Succesful = true

	return ideaFromRepo, systemMessage
}

func (ideaServise *IdeaService) DeleteIdea(id int) *entity.SystemMessage {

	systemMessage := &entity.SystemMessage{}
	erro := ideaServise.Idearepo.DeleteIdea(id)

	if erro != nil {
		systemMessage.Message = " Can't Delete The Idea "
		systemMessage.Succesful = false
		return systemMessage
	}
	systemMessage.Message = "Idea Is Deleted Succesfully "
	systemMessage.Succesful = true
	return systemMessage
}
func (ideaServise *IdeaService) VoteIdea(ideaid, voterid int) *entity.SystemMessage {
	erro := ideaServise.Idearepo.VoteIdea(ideaid, voterid)
	message := &entity.SystemMessage{}
	if erro != nil {
		message.Message = "Can't Save The Vote "
		message.Succesful = false
		return message
	}
	message.Message = "The Vote Was Succesfull "
	message.Succesful = true
	return message
}

// func (ideaServise *IdeaService) SaveCommentIdea(comment *entity.Comment) *entity.SystemMessage {

// 	systemmessage := &entity.SystemMessage{}
// 	err := ideaServise.Idearepo.SaveCommentIdea(comment)

// 	if err != nil {
// 		systemmessage.Message = "Can't Save The Message "
// 		systemmessage.Succesful = false
// 		return systemmessage
// 	}
// 	systemmessage.Message = "Succesfully Commentd "
// 	systemmessage.Succesful = true
// 	return systemmessage
// }
