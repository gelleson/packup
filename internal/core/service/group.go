package service

import (
	"github.com/gelleson/packup/internal/core/constants"
	"github.com/gelleson/packup/internal/core/dto"
	"github.com/gelleson/packup/internal/core/model"
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GroupService struct {
	db *database.Database

	createdDefaultDatabase bool
}

func NewGroupService(db *database.Database) *GroupService {
	return &GroupService{db: db}
}

func (g GroupService) Create(input dto.CreateGroupInput) (model.Group, error) {

	if !g.HasDefaultGroup() {
		return model.Group{}, errors.New("first need to create default groupObject")
	}

	if err := input.Validate(); err != nil {
		return model.Group{}, err
	}

	groupObject := model.Group{
		Name: input.Name,
	}

	if trx := g.db.Conn().Create(&groupObject); trx.Error != nil {
		return model.Group{}, trx.Error
	}

	return groupObject, nil
}

func (g GroupService) CreateDefaultGroup() (model.Group, error) {

	groupObject := model.Group{
		Model: gorm.Model{
			ID: constants.DefaultGroupId,
		},
		Name: "default",
	}

	if trx := g.db.Conn().Create(&groupObject); trx.Error != nil {
		return model.Group{}, trx.Error
	}

	return groupObject, nil
}

func (g GroupService) FindById(id uint) (model.Group, error) {

	if !g.HasDefaultGroup() {
		return model.Group{}, errors.New("first need to create default groupObject")
	}

	groupObject := model.Group{}

	if trx := g.db.Conn().First(&groupObject, "id = ?", id); trx.Error != nil {
		return model.Group{}, trx.Error
	}

	return groupObject, nil
}

func (g GroupService) Exist(id uint) bool {

	_, err := g.FindById(id)

	return err == nil
}

func (g *GroupService) HasDefaultGroup() bool {

	if g.createdDefaultDatabase {
		return true
	}

	groupObject := model.Group{}

	if trx := g.db.Conn().First(&groupObject, "id = ?", constants.DefaultGroupId); trx.Error != nil {
		return false
	}

	g.createdDefaultDatabase = true

	return true
}
