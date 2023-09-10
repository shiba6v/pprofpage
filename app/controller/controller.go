package controller

import (
	"github.com/shiba6v/pprofpage/app/repository"
)

type Controller struct {
	storage repository.ObjectStorage
}

func NewController(storage repository.ObjectStorage) Controller {
	return Controller{
		storage: storage,
	}
}
