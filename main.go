package main

//-------------------------- IMPORT --------------------------//
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
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

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/reservations", reservationHandler(db))
	http.HandleFunc("/room/add", addRoomHandler(db))
	http.HandleFunc("/reservations/add", addReservationHandler(db))
	http.HandleFunc("/room/modify", modifyReservationHandler(db))
	http.HandleFunc("/room/delete", deleteReservationHandler(db))
	http.HandleFunc("/room/list", listRoomsHandler(db))
	http.HandleFunc("/reservations_by_room", reservationsByRoomHandler(db))

	log.Println("---------------------------------------")
	log.Println("Démarrage du serveur sur le port :8095")
	log.Println("---------------------------------------")
	log.Fatal(http.ListenAndServe(":8095", nil))
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

// isRoomAvailable checks if the room is available, excluding the specified reservation ID.
func isRoomAvailable(db *sql.DB, roomID int, date, startTime, endTime string, excludeReservationID int) bool {
	var count int
	// SQL query to count conflicting reservations, excluding the current one.
	query := `SELECT COUNT(*) FROM reservations
              WHERE room_id = ? AND date = ? AND NOT (end_time <= ? OR start_time >= ?)
              AND id != ?` // Exclude the current reservation from the check.

	// Executing the query with parameters.
	err := db.QueryRow(query, roomID, date, startTime, endTime, excludeReservationID).Scan(&count)
	if err != nil {
		log.Printf("Error checking room availability: %v", err)
		return false
	}
	return count == 0 // No conflicting reservations means the room is available.
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

			// Pass 0 as excludeReservationID when adding a new reservation
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

//-------------------------- 	addRoomHandler		--------------------------//

func addRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			executeTemplate(w, "add_room.html", nil)
			return
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}
			name := r.FormValue("name")
			capacityStr := r.FormValue("capacity")
			if name == "" || capacityStr == "" {
				http.Error(w, "Le nom et la capacité de la salle sont requis.", http.StatusBadRequest)
				return
			}
			capacity, err := strconv.Atoi(capacityStr)
			if err != nil {
				http.Error(w, "Capacité invalide", http.StatusBadRequest)
				return
			}

			// Check if the room already exists
			var existingRoomId int
			checkQuery := "SELECT id FROM rooms WHERE name = ? LIMIT 1"
			err = db.QueryRow(checkQuery, name).Scan(&existingRoomId)

			if err != nil && err != sql.ErrNoRows {
				http.Error(w, "Erreur lors de la vérification de l'existence de la salle", http.StatusInternalServerError)
				return
			}

			if existingRoomId > 0 {
				http.Error(w, "Une salle avec le même nom existe déjà", http.StatusBadRequest)
				return
			}

			insertQuery := "INSERT INTO rooms (name, capacity) VALUES (?, ?)"
			_, err = db.Exec(insertQuery, name, capacity)
			if err != nil {
				http.Error(w, "Erreur lors de l'ajout de la salle", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/room/list", http.StatusSeeOther)
			return
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
		// Extract the reservation ID from the query parameters for a GET request
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
			// Handle POST request to update a reservation
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

		// Exécutez le template pour afficher les salles.
		executeTemplate(w, "list_rooms.html", data)
	}
}

//-------------------------- 	getAllRooms		--------------------------//

// getAllRooms récupère toutes les salles de la base de données.
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
		if r.Method == "GET" {
			roomID := r.URL.Query().Get("roomID")
			if roomID == "" {
				http.Error(w, "ID de salle manquant", http.StatusBadRequest)
				return
			}
			reservations, err := getReservationsByRoom(db, roomID)
			if err != nil {
				log.Printf("Error fetching reservations for room %s: %v", roomID, err)
				http.Error(w, "Erreur lors de la récupération des réservations", http.StatusInternalServerError)
				return
			}
			if len(reservations) == 0 {
				// Handle no reservations found case
				executeTemplate(w, "reservations_by_room.html", nil) // Consider how you handle no results in your template
				return
			}
			executeTemplate(w, "reservations_by_room.html", reservations)
		} else {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}
}

//-------------------------- 	geteReservationsByRoom		--------------------------//

// Fonction pour récupérer les réservations par salle depuis la base de données
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
