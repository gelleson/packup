package group

import (
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

func (g GroupService) Create(input CreateGroupInput) (Group, error) {

	if !g.HasDefaultGroup() {
		return Group{}, errors.New("first need to create default group")
	}

	if err := input.Validate(); err != nil {
		return Group{}, err
	}

	group := Group{
		Name: input.Name,
	}

	if trx := g.db.Conn().Create(&group); trx.Error != nil {
		return Group{}, trx.Error
	}

	return group, nil
}

func (g GroupService) CreateDefaultGroup() (Group, error) {

	group := Group{
		Model: gorm.Model{
			ID: DefaultGroupId,
		},
		Name: "default",
	}

	if trx := g.db.Conn().Create(&group); trx.Error != nil {
		return Group{}, trx.Error
	}

	return group, nil
}

func (g GroupService) FindById(id uint) (Group, error) {

	if !g.HasDefaultGroup() {
		return Group{}, errors.New("first need to create default group")
	}

	group := Group{}

	if trx := g.db.Conn().First(&group, "id = ?", id); trx.Error != nil {
		return Group{}, trx.Error
	}

	return group, nil
}

func (g GroupService) Exist(id uint) bool {

	_, err := g.FindById(id)

	return err == nil
}

func (g *GroupService) HasDefaultGroup() bool {

	if g.createdDefaultDatabase {
		return true
	}

	group := Group{}

	if trx := g.db.Conn().First(&group, "id = ?", DefaultGroupId); trx.Error != nil {
		return false
	}

	g.createdDefaultDatabase = true

	return true
}
