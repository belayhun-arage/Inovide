package Idea

import (
	entity "github.com/Projects/Inovide/models"
)

type IdeaRepository interface {
	CreateIdea(idea *entity.Idea) error
	DeleteIdea(id int) error
	GetIdea(id int) (*entity.Idea, error)
	VoteIdea(ideaid, voterid int) error
	SearchIdeas(text string, person *entity.Person, searchresults *[]entity.Idea) (*[]entity.Idea, error)
}
