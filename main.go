package main

//-------------------------- IMPORT --------------------------//
import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Room struct {
	ID       int
	Name     string
	Capacity int
}

type Reservation struct {
	ID        int
	RoomID    int
	Date      string
	StartTime string
	EndTime   string
}

//-------------------------- CONNEXION --------------------------//

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/projetgo")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données: %v", err)
	}
	defer db.Close()

	// Fetch reservations
	reservations, err := getAllReservations(db)
	if err != nil {
		log.Fatalf("Error fetching reservations: %v", err)
	}

	// Export to JSON
	jsonFilename := "reservations.json"
	if err = exportReservationsToJSON(reservations, jsonFilename); err != nil {
		log.Fatalf("Error exporting to JSON: %v", err)
	}

	// Export to CSV
	csvFilename := "reservations.csv"
	if err = exportReservationsToCSV(reservations, csvFilename); err != nil {
		log.Fatalf("Could not write to CSV file: %v", err)
	}

	log.Println("Reservations were successfully exported to both JSON and CSV files.")

	//---------- CSS ----------//
	staticDir := http.Dir("templates/css")
	staticHandler := http.FileServer(staticDir)
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/reservations", reservationHandler(db))
	http.HandleFunc("/room/add", addRoomHandler(db))
	http.HandleFunc("/reservations/add", addReservationHandler(db))
	http.HandleFunc("/room/modify", modifyReservationHandler(db))
	http.HandleFunc("/room/delete", deleteReservationHandler(db))
	http.HandleFunc("/room/list", listRoomsHandler(db))
	http.HandleFunc("/reservations_by_room", reservationsByRoomHandler(db))
	http.HandleFunc("/reservations_by_date", getReservationsByDateHandler(db))

	http.HandleFunc("/check_availability", checkAvailabilityHandler(db))
	log.Println("---------------------------------------")
	log.Println("Démarrage du serveur sur le port :8095")
	log.Println("---------------------------------------")
	log.Fatal(http.ListenAndServe(":8095", nil))

	// Lancer l'interface CLI

}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "home.html", nil)
}

// reservationHandler gère les réservations.
func reservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			reservations, err := getAllReservations(db)
			if err != nil {
				http.Error(w, "Erreur lors de la récupération des réservations", http.StatusInternalServerError)
				return
			}
			// Utilisez votre méthode de rendu de template ici
			executeTemplate(w, "reservations.html", reservations)

		case "POST":
			// Traitement de la création d'une nouvelle réservation
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			// Extrait les données du formulaire
			roomID, _ := strconv.Atoi(r.FormValue("roomID"))
			date := r.FormValue("date")
			startTime := r.FormValue("startTime")
			endTime := r.FormValue("endTime")

			// Vérifie la disponibilité et crée la réservation
			if isRoomAvailable(db, roomID, date, startTime, endTime, 0) {
				err := insertReservation(db, roomID, date, startTime, endTime)
				if err != nil {
					http.Error(w, "Erreur lors de la création de la réservation", http.StatusInternalServerError)
					return
				}
				// Redirection ou confirmation
				http.Redirect(w, r, "/reservations", http.StatusSeeOther)
			} else {
				http.Error(w, "La salle n'est pas disponible", http.StatusBadRequest)
			}

		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}
}

//-------------------------- IsRoomAvailable --------------------------//

func isRoomAvailable(db *sql.DB, roomID int, date, startTime, endTime string, excludeReservationID int) bool {
	var count int
	query := `SELECT COUNT(*) FROM reservations
              WHERE room_id = ? AND date = ? AND NOT (end_time <= ? OR start_time >= ?)
              AND id != ?`

	err := db.QueryRow(query, roomID, date, startTime, endTime, excludeReservationID).Scan(&count)
	if err != nil {
		log.Printf("Error checking room availability: %v", err)
		return false
	}
	return count == 0
}

//-------------------------- IsRoomAvailable --------------------------//

func insertReservation(db *sql.DB, roomID int, date, startTime, endTime string) error {
	query := `INSERT INTO reservations (room_id, date, start_time, end_time) 
              VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, roomID, date, startTime, endTime)
	if err != nil {
		log.Printf("Erreur lors de la création de la réservation : %v", err)
		return err
	}
	return nil
}

func addReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			executeTemplate(w, "add_reservation.html", nil)
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error processing form", http.StatusBadRequest)
				return
			}
			roomID, _ := strconv.Atoi(r.FormValue("roomID"))
			date := r.FormValue("date")
			startTime := r.FormValue("startTime")
			endTime := r.FormValue("endTime")

			if isRoomAvailable(db, roomID, date, startTime, endTime, 0) {
				if err := insertReservation(db, roomID, date, startTime, endTime); err != nil {
					http.Error(w, "Error creating reservation", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/reservations", http.StatusSeeOther)
			} else {
				http.Error(w, "Room is not available", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

//-------------------------- getAllReservations	 --------------------------//

func getAllReservations(db *sql.DB) ([]Reservation, error) {
	var reservations []Reservation
	query := "SELECT id, room_id, date, start_time, end_time FROM reservations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Reservation
		if err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}

//-------------------------- 	executeTemplate		--------------------------//

func executeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

//-------------------------- 	addRoomHandler		--------------------------//

func addRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			errorMessage := r.URL.Query().Get("error")
			executeTemplate(w, "add_room.html", map[string]string{"Error": errorMessage})
			return
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Redirect(w, r, "/room/add?error=Erreur lors du traitement du formulaire", http.StatusSeeOther)
				return
			}
			name := r.FormValue("name")
			capacityStr := r.FormValue("capacity")
			if name == "" || capacityStr == "" {
				http.Redirect(w, r, "/room/add?error=Le nom et la capacité de la salle sont requis.", http.StatusSeeOther)
				return
			}
			capacity, err := strconv.Atoi(capacityStr)
			if err != nil {
				http.Redirect(w, r, "/room/add?error=Capacité invalide", http.StatusSeeOther)
				return
			}

			var existingRoomId int
			checkQuery := "SELECT id FROM rooms WHERE name = ? LIMIT 1"
			err = db.QueryRow(checkQuery, name).Scan(&existingRoomId)
			if err != nil && err != sql.ErrNoRows {
				http.Redirect(w, r, "/room/add?error=Erreur lors de la vérification de l'existence de la salle", http.StatusSeeOther)
				return
			}
			if existingRoomId > 0 {
				http.Redirect(w, r, "/room/add?error=Une salle avec le même nom existe déjà", http.StatusSeeOther)
				return
			}

			insertQuery := "INSERT INTO rooms (name, capacity) VALUES (?, ?)"
			_, err = db.Exec(insertQuery, name, capacity)
			if err != nil {
				http.Redirect(w, r, "/room/add?error=Erreur lors de l'ajout de la salle", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/room/list", http.StatusSeeOther)
		} else {
			http.Error(w, "Méthode HTTP non autorisée", http.StatusMethodNotAllowed)
		}
	}
}

//-------------------------- 	deleteReservationHandler		--------------------------//

func deleteReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			reservationID := r.URL.Query().Get("id")
			query := "DELETE FROM reservations WHERE id = ?"
			_, err := db.Exec(query, reservationID)
			if err != nil {
				http.Error(w, "Erreur lors de la suppression de la réservation", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/reservations", http.StatusSeeOther)
		}
	}
}

//-------------------------- 	ModifyReservationHandler		--------------------------//

func modifyReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			reservationID := r.URL.Query().Get("id")
			if reservationID == "" {
				http.Error(w, "ID de réservation requis", http.StatusBadRequest)
				return
			}

			reservation, err := getReservationByID(db, reservationID)
			if err != nil {
				log.Printf("Erreur lors de la récupération de la réservation : %v", err)
				http.Error(w, "Erreur lors de la récupération de la réservation", http.StatusInternalServerError)
				return
			}
			executeTemplate(w, "modify_reservation.html", reservation)
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			reservationID, err := strconv.Atoi(r.FormValue("reservationID"))
			if err != nil {
				http.Error(w, "Identifiant de réservation invalide", http.StatusBadRequest)
				return
			}

			newRoomID, err := strconv.Atoi(r.FormValue("newRoomID"))
			if err != nil {
				http.Error(w, "Identifiant de salle invalide", http.StatusBadRequest)
				return
			}

			newDate := r.FormValue("newDate")
			newStartTime := r.FormValue("newStartTime")
			newEndTime := r.FormValue("newEndTime")

			if isRoomAvailable(db, newRoomID, newDate, newStartTime, newEndTime, reservationID) {
				query := `UPDATE reservations SET room_id = ?, date = ?, start_time = ?, end_time = ? WHERE id = ?`
				if _, err := db.Exec(query, newRoomID, newDate, newStartTime, newEndTime, reservationID); err != nil {
					log.Printf("Erreur lors de la modification de la réservation : %v", err)
					http.Redirect(w, r, "/reservations?message=Erreur lors de la modification de la réservation", http.StatusSeeOther)
					return
				}
				http.Redirect(w, r, "/reservations?message=Modification réussie", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/reservations?message=Modification refusée - La salle n'est pas disponible pour le créneau demandé", http.StatusSeeOther)
			}
		} else {
			http.Redirect(w, r, "Rien du tout", http.StatusSeeOther)
		}
	}
}

func getReservationByID(db *sql.DB, reservationID string) (Reservation, error) {
	var reservation Reservation
	query := `SELECT id, room_id, date, start_time, end_time FROM reservations WHERE id = ?`
	err := db.QueryRow(query, reservationID).Scan(&reservation.ID, &reservation.RoomID, &reservation.Date, &reservation.StartTime, &reservation.EndTime)
	if err != nil {
		return reservation, err
	}
	return reservation, nil
}

//-------------------------- 	ListRoomHandler		--------------------------//

func listRoomsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Méthode HTTP non autorisée", http.StatusMethodNotAllowed)
			return
		}

		rooms, err := getAllRooms(db)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des salles: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Rooms []Room
		}{
			Rooms: rooms,
		}

		executeTemplate(w, "list_rooms.html", data)
	}
}

//-------------------------- 	getAllRooms		--------------------------//

func getAllRooms(db *sql.DB) ([]Room, error) {
	var rooms []Room

	query := `SELECT id, name, capacity FROM rooms`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

//-------------------------- 	reservationsByRoomHandler		--------------------------//

func reservationsByRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		roomID := r.URL.Query().Get("roomID")
		var reservations []Reservation
		var errorMessage string

		if roomID != "" {
			var err error
			reservations, err = getReservationsByRoom(db, roomID)
			if err != nil {
				log.Printf("Error fetching reservations for room %s: %v", roomID, err)
				errorMessage = "Erreur lors de la récupération des réservations"
			} else if len(reservations) == 0 {
				errorMessage = "Aucune réservation trouvée pour cette salle."
			}
		} else {
			errorMessage = "ID de salle manquant."
		}

		// Prepare data for rendering
		data := struct {
			Reservations []Reservation
			Error        string
		}{
			Reservations: reservations,
			Error:        errorMessage,
		}

		executeTemplate(w, "reservations_by_room.html", data)
	}
}

//-------------------------- 	geteReservationsByRoom		--------------------------//

func getReservationsByRoom(db *sql.DB, roomID string) ([]Reservation, error) {
	var reservations []Reservation

	query := "SELECT id, room_id, date, start_time, end_time FROM reservations WHERE room_id = ?"

	rows, err := db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Reservation
		if err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// ------------------------------ getReservationsByDate ------------------------------//
func getReservationsByDate(db *sql.DB, date string) ([]Reservation, error) {
	var reservations []Reservation
	query := `SELECT id, room_id, date, start_time, end_time FROM reservations WHERE date = ?`
	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation Reservation
		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.Date, &reservation.StartTime, &reservation.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func getReservationsByDateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Si la méthode est GET et qu'il y a un paramètre "date", procédez à la récupération des réservations.
		if r.Method == "GET" {
			date := r.URL.Query().Get("date")
			if date == "" {
				executeTemplate(w, "reservation_by_date.html", nil)
				return
			}

			// récupérer les réservations.
			reservations, err := getReservationsByDate(db, date)
			if err != nil {
				log.Printf("Erreur lors de la récupération des réservations pour la date %s: %v", date, err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}

			// Si des réservations sont trouvées, passe  au template pour affichage.
			executeTemplate(w, "reservation_by_date.html", map[string]interface{}{
				"Date":         date,
				"Reservations": reservations,
			})
		} else {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}
}

func checkAvailabilityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		roomIDStr := r.URL.Query().Get("roomID")
		date := r.URL.Query().Get("date")
		startTime := r.URL.Query().Get("startTime")
		endTime := r.URL.Query().Get("endTime")

		var data struct {
			IsAvailable bool
			Error       string
		}

		// Check if all parameters are provided
		if roomIDStr == "" || date == "" || startTime == "" || endTime == "" {
			data.Error = "All parameters are required"
			executeTemplate(w, "availability.html", data)
			return
		}

		// Attempt to convert roomID to integer
		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			data.Error = "Invalid room ID"
			executeTemplate(w, "availability.html", data)
			return
		}

		// Check room availability and set the data
		data.IsAvailable = isRoomAvailable(db, roomID, date, startTime, endTime, 0)
		executeTemplate(w, "availability.html", data)
	}
}

//------------------------------	Export JSON		------------------------------------//

func exportReservationsToJSON(reservations []Reservation, filename string) error {
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

// ------------------------------ 	Export CSV		------------------------------------//

func exportReservationsToCSV(reservations []Reservation, filename string) error {
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
