package migration

import (
	"fmt"
	"medicalHistory/database"
	"medicalHistory/model/entity"

	"log"
)

func Migration() {
	//database.DB merupakan variable yang di assign di database.go
	err := database.DB.AutoMigrate(

		&entity.MedicalHistory{},
	)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrated")
}
