package exportlogic

import (
	"Reserve-Go/menulogic"
	"Reserve-Go/reservationlogic"
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func ExportReservationsAsCSV(db *sql.DB, filename string, scanner *bufio.Scanner) error {
	reservations, err := reservationlogic.GetAllReservations(db)
	if err != nil {
		log.Printf("Error fetching reservations: %v", err)
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating CSV file: %v", err)
		return err
	}
	defer func(file *os.File) {
		expErr := file.Close()
		if expErr != nil {
			log.Printf("Erreur: %v", expErr)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "RoomID", "Date", "StartTime", "EndTime"}
	if err := writer.Write(header); err != nil {
		log.Printf("Error writing header to CSV: %v", err)
		return err
	}

	for _, reservation := range reservations {
		record := []string{
			strconv.Itoa(reservation.ID),
			strconv.Itoa(reservation.RoomID),
			reservation.Date,
			reservation.StartTime,
			reservation.EndTime,
		}
		if err := writer.Write(record); err != nil {
			log.Printf("Error writing record to CSV: %v", err)
			return err
		}
	}

	log.Printf("Reservations successfully exported to %s", filename)
	menulogic.NavigationOptions(db, scanner)
	return nil
}

func ExportReservationsAsJSON(db *sql.DB, filename string, scanner *bufio.Scanner) error {
	reservations, err := reservationlogic.GetAllReservations(db)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(reservations, "", "    ")
	if err != nil {
		return err
	}
	menulogic.NavigationOptions(db, scanner)
	return ioutil.WriteFile(filename, data, 0644)
}
