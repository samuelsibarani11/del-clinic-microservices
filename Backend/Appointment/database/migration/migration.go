package migration

import (
	"appointment/database"
	"appointment/model/entity"
	"fmt"
	"log"
)

func Migration() {
	//database.DB merupakan variable yang di assign di database.go
	err := database.DB.AutoMigrate(
		&entity.Appointment{},
	)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrated")
}
