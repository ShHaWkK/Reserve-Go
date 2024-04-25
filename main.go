package main

//-------------------------- IMPORT --------------------------//
import (
	"Reserve-Go/dtb"
	"Reserve-Go/exportlogic"
	"Reserve-Go/reservationlogic"
	"Reserve-Go/roomlogic"
	"Reserve-Go/utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	db, err := dtb.ConnectToDB()
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données: %v", err)
	}
	defer db.Close()

	// Fetch reservations
	reservations, err := reservationlogic.GetAllReservations(db)
	if err != nil {
		log.Fatalf("Error fetching reservations: %v", err)
	}

	// Export to JSON
	jsonFilename := "reservations.json"
	if err = exportlogic.ExportReservationsToJSON(reservations, jsonFilename); err != nil {
		log.Fatalf("Error exporting to JSON: %v", err)
	}

	// Export to CSV
	csvFilename := "reservations.csv"
	if err = exportlogic.ExportReservationsToCSV(reservations, csvFilename); err != nil {
		log.Fatalf("Could not write to CSV file: %v", err)
	}

	log.Println("Reservations were successfully exported to both JSON and CSV files.")

	//---------- CSS ----------//
	staticDir := http.Dir("templates/css")
	staticHandler := http.FileServer(staticDir)
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	http.HandleFunc("/", utils.HomeHandler)
	http.HandleFunc("/reservations", reservationlogic.ReservationHandler(db))
	http.HandleFunc("/room/add", roomlogic.AddRoomHandler(db))
	http.HandleFunc("/reservations/add", reservationlogic.AddReservationHandler(db))
	http.HandleFunc("/room/modify", reservationlogic.ModifyReservationHandler(db))
	http.HandleFunc("/room/delete", reservationlogic.DeleteReservationHandler(db))
	http.HandleFunc("/room/list", roomlogic.ListRoomsHandler(db))
	http.HandleFunc("/reservations_by_room", reservationlogic.ReservationsByRoomHandler(db))
	http.HandleFunc("/reservations_by_date", reservationlogic.GetReservationsByDateHandler(db))
	http.HandleFunc("/check_availability", roomlogic.CheckAvailabilityHandler(db))
	http.HandleFunc("/download/csv", exportlogic.DownloadCSVHandler(db))
	http.HandleFunc("/download/json", exportlogic.DownloadJSONHandler(db))
	http.HandleFunc("/download", utils.DownloadHandler)

	log.Println("---------------------------------------")
	log.Println("Démarrage du serveur sur le port :8095")
	log.Println("---------------------------------------")
	log.Fatal(http.ListenAndServe(":8095", nil))

}
