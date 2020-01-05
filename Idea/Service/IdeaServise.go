package service

import (
	repository "github.com/Samuael/Projects/Inovide/Idea/Repository"
	//entity "github.com/Samuael/Projects/Inovide/models"
)

type IdeaService struct {
	idearepo *repository.IdeaRepo
}

func NewIdeaService(idearep *repository.IdeaRepo) *IdeaService {
	return &IdeaService{idearepo: idearep}
}
