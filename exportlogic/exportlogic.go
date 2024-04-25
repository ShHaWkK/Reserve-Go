package exportlogic

import (
	"Reserve-Go/models"
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ExportReservationsToJSON(reservations []models.Reservation, filename string) error {
	jsonData, err := json.MarshalIndent(reservations, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling reservations: %v", err)
	}

	if err = ioutil.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	log.Printf("Reservations exported to %s successfully.\n", filename)
	return nil
}

func ExportReservationsToCSV(reservations []models.Reservation, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "RoomID", "Date", "StartTime", "EndTime"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, reservation := range reservations {
		record := []string{
			fmt.Sprint(reservation.ID),
			fmt.Sprint(reservation.RoomID),
			reservation.Date,
			reservation.StartTime,
			reservation.EndTime,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	log.Printf("Reservations exported to %s successfully.\n", filename)
	return nil
}

func generateCSVData(db *sql.DB) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{"ID", "RoomID", "Date", "StartTime", "EndTime"})
	rows, err := db.Query("SELECT id, room_id, date, start_time, end_time FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime)
		if err != nil {
			return nil, err
		}
		writer.Write([]string{
			strconv.Itoa(r.ID),
			strconv.Itoa(r.RoomID),
			r.Date,
			r.StartTime,
			r.EndTime,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func generateJSONData(db *sql.DB) ([]byte, error) {
	var reservations []models.Reservation
	rows, err := db.Query("SELECT id, room_id, date, start_time, end_time FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var res models.Reservation
		if err := rows.Scan(&res.ID, &res.RoomID, &res.Date, &res.StartTime, &res.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(reservations)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func DownloadCSVHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("test")
		data, err := generateCSVData(db)
		if err != nil {
			http.Error(w, "Failed to generate CSV", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=reservations_web_download.csv")
		w.Write(data)
	}
}

func DownloadJSONHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := generateJSONData(db)
		if err != nil {
			http.Error(w, "Failed to generate JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment;filename=Export.json")
		w.Write(data)
	}
}
