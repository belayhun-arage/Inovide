package ideaService

import (
	IdeaRepository "github.com/Samuael/Projects/Inovide/Idea/Repository"
	entity "github.com/Samuael/Projects/Inovide/models"
)

type IdeaService struct {
	idearepo *IdeaRepository.IdeaRepo
}

func NewIdeaService(idearep *IdeaRepository.IdeaRepo) *IdeaService {
	return &IdeaService{idearepo: idearep}
}

func (ideaServise *IdeaService) CreateIdea(idea *entity.Idea) *entity.SystemMessage {
	var message = entity.SystemMessage{}
	if idea.Title == "" {
		message.Message = "The Input Is Not Fully FIlled Please Submitt Again Filling The Data Appropriately"
		message.Succesful = false
	} else {
		er := ideaServise.idearepo.CreateIdea(idea)
		if er != nil {
			panic(er)
		}
	}
	message.Message = "Succesfully Inserted "
	message.Succesful = true
	return &message
}
