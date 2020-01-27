package Idea

import (
	entity "github.com/Projects/Inovide/models"
)

type IdeaService interface {
	CreateIdea(idea *entity.Idea) *entity.SystemMessage
	GetIdea(idea *entity.Idea, id int) *entity.SystemMessage
	DeleteIdea(id int) *entity.SystemMessage
	VoteIdea(ideaid, voterid int) *entity.SystemMessage
	SearchResult(searchingtext string, person *entity.Person, searchresults *[]entity.Idea) *entity.SystemMessage
	MyIdeas(ideaowner int, ideas *[]entity.Idea) *entity.SystemMessage
}
