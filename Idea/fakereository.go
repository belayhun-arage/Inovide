package Idea

import (
	entity "github.com/Projects/Inovide/models"
)

type FakeIdeaRepository interface {
	CreateFakeIdea(idea *entity.Idea) error
}
