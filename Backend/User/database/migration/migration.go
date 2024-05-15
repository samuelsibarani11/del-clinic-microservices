package migration

import (
	"fmt"
	"log"
	"user/database"
	"user/models/entity"
)

func Migration() {

	err := database.DB.AutoMigrate(&entity.User{})

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database Migrated")
}
