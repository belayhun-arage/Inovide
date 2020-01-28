package ideaService

import (
	"github.com/Projects/Inovide/Idea"
)

type FakeIdeaService struct {
	FakeIdearepo Idea.IdeaRepository
}

func NewFakeIdeaService(fakeidearep Idea.IdeaRepository) *FakeIdeaService {
	return &FakeIdeaService{FakeIdearepo: fakeidearep}
}
