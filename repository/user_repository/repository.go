package user_repository

import (
	"github.com/nabilsea/hacktiv8-final-project/entity"
	"github.com/nabilsea/hacktiv8-final-project/pkg/errs"
)

type UserRepository interface {
	Login(user *entity.User) (*entity.User, errs.MessageErr)
	Register(user *entity.User) (*entity.User, errs.MessageErr)
	GetUserByIDAndEmail(user *entity.User) (*entity.User, errs.MessageErr)
	UpdateUserData(userId uint, user *entity.User) (*entity.User, errs.MessageErr)
	DeleteUser(userId uint) errs.MessageErr
}
