package migration

import (
	"fmt"
	"log"
	"medicine/database"
	"medicine/model/entity"
)

func Migration() {
	//database.DB merupakan variable yang di assign di database.go
	err := database.DB.AutoMigrate(
		&entity.Medicine{},
	)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrated")
}
