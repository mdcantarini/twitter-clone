package repository

import (
	"github.com/mdcantarini/twitter-clone/internal/user"
)

type Repository interface {
	InsertUser(user *user.User) (*user.User, error)
	GetUser(id uint) (*user.User, error)
}
