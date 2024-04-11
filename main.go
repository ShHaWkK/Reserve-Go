package main

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

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database_name")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données:", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		executeTemplate(w, "home.html", nil)
	})

	http.HandleFunc("/reservations", reservationHandler(db))
	http.HandleFunc("/reservations/add-room", addRoomHandler(db))
	http.HandleFunc("/modify-reservation", modifyReservationHandler(db))
	http.HandleFunc("/delete-reservation", deleteReservationHandler(db))

	// Démarrage du serveur
	log.Println("Démarrage du serveur sur :8095")
	log.Fatal(http.ListenAndServe(":8095", nil))
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
			if isRoomAvailable(db, roomID, date, startTime, endTime) {
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
func isRoomAvailable(db *sql.DB, roomID int, date, startTime, endTime string) bool {
	var count int
	query := `SELECT COUNT(*) FROM reservations 
              WHERE room_id = ? 
              AND date = ? 
              AND NOT (start_time >= ? OR end_time <= ?)`

	err := db.QueryRow(query, roomID, date, endTime, startTime).Scan(&count)
	if err != nil {
		log.Printf("Erreur lors de la vérification de la disponibilité : %v", err)
		return false
	}
	return count == 0
}
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

func executeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func addRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}
			name := r.FormValue("name")
			capacity, err := strconv.Atoi(r.FormValue("capacity"))
			if err != nil {
				http.Error(w, "Capacité invalide", http.StatusBadRequest)
				return
			}
			query := "INSERT INTO rooms (name, capacity) VALUES (?, ?)"
			_, err = db.Exec(query, name, capacity)
			if err != nil {
				http.Error(w, "Erreur lors de l'ajout de la salle", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			executeTemplate(w, "add_room.html", nil)
		}
	}
}
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
func modifyReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Traiter le formulaire de modification
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

			// Vérifiez si la nouvelle salle est disponible pour le nouveau créneau
			if isRoomAvailable(db, newRoomID, newDate, newStartTime, newEndTime) {
				query := `UPDATE reservations SET room_id = ?, date = ?, start_time = ?, end_time = ? WHERE id = ?`
				_, err = db.Exec(query, newRoomID, newDate, newStartTime, newEndTime, reservationID)
				if err != nil {
					http.Error(w, "Erreur lors de la modification de la réservation", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/reservations", http.StatusSeeOther)
			} else {
				http.Error(w, "La salle n'est pas disponible pour le créneau demandé", http.StatusBadRequest)
			}

		} else {
			// Afficher le formulaire de modification (pour simplifier, on ne le montre pas ici)
		}
	}
}
