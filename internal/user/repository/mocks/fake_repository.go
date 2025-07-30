package mocks

import (
	"github.com/mdcantarini/twitter-clone/internal/user/model"
)

type FakeSqlRepository struct {
	InsertUserFunc func(user *model.User) (*model.User, error)
	GetUserFunc    func(id uint) (*model.User, error)
}

func (f FakeSqlRepository) InsertUser(user *model.User) (*model.User, error) {
	if f.InsertUserFunc != nil {
		return f.InsertUserFunc(user)
	}
	user.ID = 1
	return user, nil
}

func (f FakeSqlRepository) GetUser(id uint) (*model.User, error) {
	if f.GetUserFunc != nil {
		return f.GetUserFunc(id)
	}
	return &model.User{
		ID:          id,
		Username:    "testuser",
		DisplayName: "Test User",
	}, nil
}
