package reservationlogic

import (
	"Reserve-Go/models"
	"Reserve-Go/roomlogic"
	"Reserve-Go/utils"
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

// reservationHandler gère les réservations.
func ReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			reservations, err := GetAllReservations(db)
			if err != nil {
				http.Error(w, "Erreur lors de la récupération des réservations", http.StatusInternalServerError)
				return
			}
			// Utilisez votre méthode de rendu de template ici
			utils.ExecuteTemplate(w, "reservations.html", reservations)

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
			if roomlogic.IsRoomAvailable(db, roomID, date, startTime, endTime, 0) {
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

func AddReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			utils.ExecuteTemplate(w, "add_reservation.html", nil)
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error processing form", http.StatusBadRequest)
				return
			}
			roomID, _ := strconv.Atoi(r.FormValue("roomID"))
			date := r.FormValue("date")
			startTime := r.FormValue("startTime")
			endTime := r.FormValue("endTime")

			if roomlogic.IsRoomAvailable(db, roomID, date, startTime, endTime, 0) {
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

func GetAllReservations(db *sql.DB) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := "SELECT id, room_id, date, start_time, end_time FROM reservations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Reservation
		if err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}

func DeleteReservationHandler(db *sql.DB) http.HandlerFunc {
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

func ModifyReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			reservationID := r.URL.Query().Get("id")
			if reservationID == "" {
				http.Error(w, "ID de réservation requis", http.StatusBadRequest)
				return
			}

			reservation, err := GetReservationByID(db, reservationID)
			if err != nil {
				log.Printf("Erreur lors de la récupération de la réservation : %v", err)
				http.Error(w, "Erreur lors de la récupération de la réservation", http.StatusInternalServerError)
				return
			}
			utils.ExecuteTemplate(w, "modify_reservation.html", reservation)
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

			if roomlogic.IsRoomAvailable(db, newRoomID, newDate, newStartTime, newEndTime, reservationID) {
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

func GetReservationByID(db *sql.DB, reservationID string) (models.Reservation, error) {
	var reservation models.Reservation
	query := `SELECT id, room_id, date, start_time, end_time FROM reservations WHERE id = ?`
	err := db.QueryRow(query, reservationID).Scan(&reservation.ID, &reservation.RoomID, &reservation.Date, &reservation.StartTime, &reservation.EndTime)
	if err != nil {
		return reservation, err
	}
	return reservation, nil
}

func ReservationsByRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		roomID := r.URL.Query().Get("roomID")
		var reservations []models.Reservation
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
			Reservations []models.Reservation
			Error        string
		}{
			Reservations: reservations,
			Error:        errorMessage,
		}

		utils.ExecuteTemplate(w, "reservations_by_room.html", data)
	}
}

func getReservationsByRoom(db *sql.DB, roomID string) ([]models.Reservation, error) {
	var reservations []models.Reservation

	query := "SELECT id, room_id, date, start_time, end_time FROM reservations WHERE room_id = ?"

	rows, err := db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Reservation
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

func getReservationsByDate(db *sql.DB, date string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := `SELECT id, room_id, date, start_time, end_time FROM reservations WHERE date = ?`
	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.Date, &reservation.StartTime, &reservation.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func GetReservationsByDateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Si la méthode est GET et qu'il y a un paramètre "date", procédez à la récupération des réservations.
		if r.Method == "GET" {
			date := r.URL.Query().Get("date")
			if date == "" {
				utils.ExecuteTemplate(w, "reservation_by_date.html", nil)
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
			utils.ExecuteTemplate(w, "reservation_by_date.html", map[string]interface{}{
				"Date":         date,
				"Reservations": reservations,
			})
		} else {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}
}
