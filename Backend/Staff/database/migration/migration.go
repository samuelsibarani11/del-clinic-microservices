package migration

import (
	"fmt"
	"log"
	"staff/database"
	"staff/models/entity"
)

func Migration() {

	err := database.DB.AutoMigrate(&entity.Staff{})

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database Migrated")
}
