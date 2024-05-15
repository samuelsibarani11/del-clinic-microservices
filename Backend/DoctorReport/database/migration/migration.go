package migration

import (
	"doctorReport/database"
	"doctorReport/model/entity"
	"fmt"
	"log"
)

func Migration() {
	//database.DB merupakan variable yang di assign di database.go
	err := database.DB.AutoMigrate(

		&entity.DoctorReport{},
	)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Database Migrated")
}
