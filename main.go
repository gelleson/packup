package main

import (
	"fmt"
	"github.com/gelleson/packup/internal/acl"
	"github.com/gelleson/packup/internal/group"
	"github.com/gelleson/packup/internal/user"
	"github.com/gelleson/packup/pkg/database"
	"gorm.io/gorm"
	"log"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db := database.NewDatabase(database.Config{
		DSN: "test2.db",
	})

	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}

	gs := group.NewService(db)
	service := user.NewService(db, gs)
	c := acl.NewService(db, gs)

	//rule, _ := c.Create(acl.Rule{
	//	Operation: acl.ReadOps,
	//	Resource: "t",
	//	GroupID: group.DefaultGroupId,
	//})
	//

	g, err := service.FindById(3)

	if err == gorm.ErrRecordNotFound {
		fmt.Println(2331)
		return
	}

	if c.HasDefaultRules() {
		fmt.Println(1)
	}

	fmt.Println(g.Group.Name, err, c.Can(1, acl.ReadOps, "acl.GroupResource"))
}
