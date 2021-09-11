package service

import (
	"github.com/gelleson/packup/internal/core/constants"
	"github.com/gelleson/packup/internal/core/dto"
	"github.com/gelleson/packup/internal/core/model"
	"github.com/gelleson/packup/pkg/database"
	"time"
)

type groupUserService interface {
	Exist(id uint) bool
}

type UserService struct {
	db           *database.Database
	groupService groupUserService
}

func NewUserService(db *database.Database, groupService groupUserService) *UserService {
	return &UserService{db: db, groupService: groupService}
}

func (u UserService) Create(input dto.CreateUserInput) (model.User, error) {

	if err := input.Validate(); err != nil {
		return model.User{}, err
	}

	user := model.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if input.HasGroup() && u.groupService.Exist(input.GroupId) {
		user.GroupID = input.GroupId
	} else {
		user.GroupID = constants.DefaultGroupId
	}

	if trx := u.db.Conn().Create(&user); trx.Error != nil {
		return model.User{}, trx.Error
	}

	return user, nil
}

func (u UserService) FindById(userId uint) (model.User, error) {

	user := model.User{}

	if trx := u.db.Conn().Preload("Group").First(&user, "id = ?", userId); trx.Error != nil {
		return model.User{}, trx.Error
	}

	return user, nil
}

func (u UserService) FindByEmail(email string) (model.User, error) {

	user := model.User{}

	if trx := u.db.Conn().Preload("Group").First(&user, "email = ?", email); trx.Error != nil {
		return model.User{}, trx.Error
	}

	return user, nil
}

func (u UserService) SetLoggedTime(userId uint, t time.Time) error {

	if trx := u.db.Conn().Model(&model.User{}).Where("id = ?", userId).Update("last_logged", t); trx.Error != nil {
		return trx.Error
	}

	return nil
}
