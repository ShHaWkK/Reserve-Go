package roomlogic

import (
	"Reserve-Go/models"
	"Reserve-Go/utils"
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func IsRoomAvailable(db *sql.DB, roomID int, date, startTime, endTime string, excludeReservationID int) bool {
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

func AddRoomHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			errorMessage := r.URL.Query().Get("error")
			utils.ExecuteTemplate(w, "add_room.html", map[string]string{"Error": errorMessage})
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

func ListRoomsHandler(db *sql.DB) http.HandlerFunc {
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
			Rooms []models.Room
		}{
			Rooms: rooms,
		}

		utils.ExecuteTemplate(w, "list_rooms.html", data)
	}
}

func getAllRooms(db *sql.DB) ([]models.Room, error) {
	var rooms []models.Room

	query := `SELECT id, name, capacity FROM rooms`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
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

func CheckAvailabilityHandler(db *sql.DB) http.HandlerFunc {
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
		if roomIDStr == "" || date == "" || startTime == "" || endTime == "" {
			data.Error = "All parameters are required"
			utils.ExecuteTemplate(w, "availability.html", data)
			return
		}

		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			data.Error = "Invalid room ID"
			utils.ExecuteTemplate(w, "availability.html", data)
			return
		}

		data.IsAvailable = IsRoomAvailable(db, roomID, date, startTime, endTime, 0)
		utils.ExecuteTemplate(w, "availability.html", data)
	}
}
