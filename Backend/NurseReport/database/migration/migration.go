package migration

import (
	"fmt"
	"log"
	"nurseReport/database"
	"nurseReport/model/entity"
)

func Migration() {
	//database.DB merupakan variable yang di assign di database.go
	err := database.DB.AutoMigrate(
		&entity.NurseReport{},
	)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrated")
}
