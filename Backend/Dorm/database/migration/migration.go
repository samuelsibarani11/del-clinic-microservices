package migration

import (
	"dorm/Models/entity"
	"dorm/database"
	"fmt"
	"log"
)

func Migration() {
	//database.DB merupakan variable yang di assign di database.go
	err := database.DB.AutoMigrate(
		&entity.Dorm{},
	)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrated")
}
