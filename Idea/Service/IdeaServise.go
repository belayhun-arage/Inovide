package ideaService

import (
	"github.com/Projects/Inovide/Idea"

	entity "github.com/Projects/Inovide/models"
)

type IdeaService struct {
	Idearepo Idea.IdeaRepository
}

func NewIdeaService(idearep Idea.IdeaRepository) *IdeaService {
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

func (ideaServise *IdeaService) GetIdea(idea *entity.Idea, id int) *entity.SystemMessage {
	systemMessage := &entity.SystemMessage{}

	numberofaffected := ideaServise.Idearepo.GetIdea(idea)

	if numberofaffected <= 0 {
		systemMessage.Message = "No Can't Get  the Idea "
		systemMessage.Succesful = false
		return systemMessage
	}
	systemMessage.Message = "Idea Succesfully Found"
	systemMessage.Succesful = true
	return systemMessage
}

func (ideaServise *IdeaService) DeleteIdea(idea *entity.Idea) *entity.SystemMessage {

	systemMessage := &entity.SystemMessage{}
	RowsAffected := ideaServise.Idearepo.DeleteIdea(idea)

	if RowsAffected <= 0 {
		systemMessage.Message = " Can't Delete The Idea "
		systemMessage.Succesful = false
		return systemMessage
	}
	systemMessage.Message = "Idea Is Deleted Succesfully "
	systemMessage.Succesful = true
	return systemMessage
}

func (ideaservice *IdeaService) UpdateIdea(idea *entity.Idea) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}
	val := ideaservice.Idearepo.UpdateIdea(idea)
	if val <= 0 {
		systemmessage.Message = "Can't Update The Idea "
		systemmessage.Succesful = false
		// }  else if val == 2 {

	} else { // 	systemmessage.Message = "There Is No Idea By This id Owned By You"
		// 	systemmessage.Succesful = false
		// }
		systemmessage.Message = "SuccesFully Updated"
		systemmessage.Succesful = true

	}
	return systemmessage
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
func (ideaservice *IdeaService) SearchResult(searchingtext string, person *entity.Person, searchresults *[]entity.Idea) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}
	_, erro := ideaservice.Idearepo.SearchIdeas(searchingtext, person, searchresults)
	if erro != nil {
		systemmessage.Succesful = false
		return systemmessage
	}
	systemmessage.Succesful = true
	return systemmessage
}
