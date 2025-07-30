package repository

import (
	"github.com/mdcantarini/twitter-clone/internal/user/model"
)

type Repository interface {
	InsertUser(user *model.User) (*model.User, error)
	GetUser(id uint) (*model.User, error)
}
