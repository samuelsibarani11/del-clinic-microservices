package migration

import (
	"fmt"
	"reminder/database"
	"reminder/model/entity"

	"log"
)

func Migration() {
	//database.DB merupakan variable yang di assign di database.go
	err := database.DB.AutoMigrate(
		&entity.Reminder{},
	)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrated")
}
