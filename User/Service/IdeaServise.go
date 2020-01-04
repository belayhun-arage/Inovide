package service

import (
	repository "github.com/BELAY-hun/Projects/Inovide/User/Repository"
	//entity "github.com/BELAY-hun/Projects/Inovide/models"
)

type IdeaService struct {
	idearepo *repository.IdeaRepo
}

func NewIdeaService(idearep *repository.IdeaRepo) *IdeaService {
	return &IdeaService{idearepo: idearep}
}
