package user

import (
	"github.com/gelleson/packup/internal/modules/group"
	"github.com/gelleson/packup/pkg/database"
	"time"
)

type groupService interface {
	Exist(id uint) bool
}

type UserService struct {
	db           *database.Database
	groupService groupService
}

func NewUserService(db *database.Database, groupService groupService) *UserService {
	return &UserService{db: db, groupService: groupService}
}

func (u UserService) Create(input CreateUserInput) (User, error) {

	if err := input.Validate(); err != nil {
		return User{}, err
	}

	user := User{
		Email:    input.Email,
		Password: input.Password,
	}

	if input.HasGroup() && u.groupService.Exist(input.GroupId) {
		user.GroupID = input.GroupId
	} else {
		user.GroupID = group.DefaultGroupId
	}

	if trx := u.db.Conn().Create(&user); trx.Error != nil {
		return User{}, trx.Error
	}

	return user, nil
}

func (u UserService) FindById(userId uint) (User, error) {

	user := User{}

	if trx := u.db.Conn().Preload("Group").First(&user, "id = ?", userId); trx.Error != nil {
		return User{}, trx.Error
	}

	return user, nil
}

func (u UserService) FindByEmail(email string) (User, error) {

	user := User{}

	if trx := u.db.Conn().Preload("Group").First(&user, "email = ?", email); trx.Error != nil {
		return User{}, trx.Error
	}

	return user, nil
}

func (u UserService) SetLoggedTime(userId uint, t time.Time) error {

	if trx := u.db.Conn().Model(&User{}).Where("id = ?", userId).Update("last_logged", t); trx.Error != nil {
		return trx.Error
	}

	return nil
}
