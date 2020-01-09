package ideaService

import (
	"fmt"

	IdeaRepository "github.com/Projects/Inovide/Idea/Repository"
	entity "github.com/Projects/Inovide/models"
)

type IdeaService struct {
	idearepo *IdeaRepository.PsqlIdeaRepository
}

func NewIdeaService(idearep *IdeaRepository.PsqlIdeaRepository) *IdeaService {
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
			fmt.Println("error in ideaservise")
			panic(er)
		}
	}
	message.Message = " Your idea is succesfully posted "
	message.Succesful = true
	return &message
}
