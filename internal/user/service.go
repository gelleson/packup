package user

import (
	"github.com/gelleson/packup/internal/group"
	"github.com/gelleson/packup/pkg/database"
	"time"
)

type groupService interface {
	Exist(id uint) bool
}

type Service struct {
	db           *database.Database
	groupService groupService
}

func NewService(db *database.Database, groupService groupService) *Service {
	return &Service{db: db, groupService: groupService}
}

func (s Service) Create(input CreateUserInput) (User, error) {

	if err := input.Validate(); err != nil {
		return User{}, err
	}

	user := User{
		Email:    input.Email,
		Password: input.Password,
	}

	if input.HasGroup() && s.groupService.Exist(input.GroupId) {
		user.GroupID = input.GroupId
	} else {
		user.GroupID = group.DefaultGroupId
	}

	if trx := s.db.Conn().Create(&user); trx.Error != nil {
		return User{}, trx.Error
	}

	return user, nil
}

func (s Service) FindById(userId uint) (User, error) {

	user := User{}

	if trx := s.db.Conn().Preload("Group").First(&user, "id = ?", userId); trx.Error != nil {
		return User{}, trx.Error
	}

	return user, nil
}

func (s Service) FindByEmail(email string) (User, error) {

	user := User{}

	if trx := s.db.Conn().First(&user, "email = ?", email); trx.Error != nil {
		return User{}, trx.Error
	}

	return user, nil
}

func (s Service) SetLoggedTime(userId uint, t time.Time) error {

	if trx := s.db.Conn().Model(&User{}).Where("id = ?", userId).Update("last_logged", t); trx.Error != nil {
		return trx.Error
	}

	return nil
}
