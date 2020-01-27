package Idea

import (
	entity "github.com/Projects/Inovide/models"
)

type IdeaRepository interface {
	CreateIdea(idea *entity.Idea) error
	DeleteIdea(idea *entity.Idea) int64
	GetIdea(idea *entity.Idea) int64
	VoteIdea(ideaid, voterid int) error
	SearchIdeas(text string, person *entity.Person, searchresults *[]entity.Idea) (*[]entity.Idea, error)
	UpdateIdea(idea *entity.Idea) int64
	MyIdeas(ideaowner int, ideas *[]entity.Idea) int64
}
