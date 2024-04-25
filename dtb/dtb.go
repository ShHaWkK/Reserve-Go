package dtb

import (
	"Reserve-Go/utils"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func ConnectToDB() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		"user",
		"password",
		"localhost",
		"3306",
		"projetgo",
	)

	log.Printf("Connecting to database with connection string: %s", connectionString)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			utils.ColorLog(utils.ColorGreen, "Successfully connected to the database.")
			return db, nil
		}
		utils.ColorLog(utils.ColorRed, "Failed to connect to the database, retrying in 1 second...")
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("error verifying connection to the database: %v", err)
}
