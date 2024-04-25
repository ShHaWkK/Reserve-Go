package reservationlogic

import (
	"Reserve-Go/menulogic"
	"Reserve-Go/models"
	"Reserve-Go/roomlogic"
	"Reserve-Go/utils"
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func CreateReservation(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println(utils.ColorString(utils.ColorBlue, strings.Repeat("-", 35)))
	fmt.Println("Création d'une réservation...")
	fmt.Println(utils.ColorString(utils.ColorBlue, strings.Repeat("-", 35)))

	fmt.Println("Entrez l'ID de la salle :")
	scanner.Scan()
	roomID := scanner.Text()

	fmt.Println("Entrez la date de réservation (YYYY-MM-DD) :")
	scanner.Scan()
	date := scanner.Text()

	fmt.Println("Entrez l'heure de début (HH:MM:SS) :")
	scanner.Scan()
	startTime := scanner.Text()

	fmt.Println("Entrez l'heure de fin (HH:MM:SS) :")
	scanner.Scan()
	endTime := scanner.Text()

	if roomlogic.IsRoomAvailable(db, roomID, date, startTime, endTime) {
		InsertReservation(db, roomID, date, startTime, endTime)
		fmt.Println("Réservation créée avec succès.")
	} else {
		fmt.Println("La salle n'est pas disponible pour le créneau demandé.")
	}
	menulogic.NavigationOptions(db, scanner)
}

func ViewReservationsByRoom(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'ID de la salle (nombre entier) : ")
	scanner.Scan()
	roomIDStr := scanner.Text()

	// Conversion de l'ID de la salle en int
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		fmt.Println("Erreur : ID de salle invalide. Veuillez entrer un nombre.")
		return
	}

	// Vérifier si la salle existe
	if !roomlogic.IsRoomExists(db, roomID) {
		fmt.Println("La salle avec l'ID", roomID, "n'existe pas.")
		return
	}

	// Appel à getReservationsByRoom avec l'ID de la salle
	reservations, err := GetReservationsByRoom(db, roomID)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des réservations :", err)
		return
	}

	// Affichage des réservations
	if len(reservations) == 0 {
		fmt.Println("Aucune réservation trouvée pour la salle", roomID)
		return
	}

	fmt.Println("Réservations pour la salle", roomID)
	for _, reservation := range reservations {
		fmt.Printf("ID: %d, Date: %s, Début: %s, Fin: %s\n", reservation.ID, reservation.Date, reservation.StartTime, reservation.EndTime)
	}
}

func GetReservationsByRoom(db *sql.DB, roomID int) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := "SELECT id, room_id, date, start_time, end_time FROM reservations WHERE room_id = ?"
	rows, err := db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		rowErr := rows.Close()
		if rowErr != nil {
			log.Printf("Erreur: %v", rowErr)
		}
	}(rows)

	for rows.Next() {
		var r models.Reservation
		if err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}

func GetReservationsByDate(db *sql.DB, date string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := "SELECT id, room_id, date, start_time, end_time FROM reservations WHERE date = ?"
	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		rowErr := rows.Close()
		if rowErr != nil {
			log.Printf("Erreur: %v", rowErr)
		}
	}(rows)

	for rows.Next() {
		var r models.Reservation
		if err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}

func InsertReservation(db *sql.DB, roomID, date, startTime, endTime string) {
	query := `INSERT INTO reservations (room_id, date, start_time, end_time) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, roomID, date, startTime, endTime)
	if err != nil {
		log.Printf("Erreur lors de la création de la réservation : %v", err)
	} else {

	}

}

func CancelReservation(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'identifiant de la réservation à annuler : ")
	scanner.Scan()
	reservationID := scanner.Text()

	// Vérification de l'existence de la réservation avant de tenter de l'annuler
	if ReservationExists(db, reservationID) {
		err := DeleteReservation(db, reservationID)
		if err != nil {
			fmt.Println("Erreur lors de l'annulation de la réservation :", err)
		} else {
			fmt.Println("Réservation annulée avec succès.")
		}
	} else {
		fmt.Println("Aucune réservation trouvée avec cet identifiant.")
	}
	menulogic.NavigationOptions(db, scanner)
}

func DeleteReservation(db *sql.DB, reservationID string) error {
	query := "DELETE FROM reservations WHERE id = ?"
	_, err := db.Exec(query, reservationID)
	return err
}

func ViewReservations(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Visualisation des réservations:")

	query := `SELECT r.id, r.room_id, r.date, r.start_time, r.end_time 
              FROM reservations r
              ORDER BY r.date, r.start_time`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Erreur lors de la récupération des réservations : %v", err)
		return
	}
	defer func(rows *sql.Rows) {
		rowErr := rows.Close()
		if rowErr != nil {
			log.Printf("Erreur: %v", rowErr)
		}
	}(rows)

	for rows.Next() {
		var reservation models.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.Date, &reservation.StartTime, &reservation.EndTime); err != nil {
			log.Printf("Erreur lors de la lecture des données de la réservation : %v", err)
			continue
		}
		fmt.Printf("ID: %d, Salle: %d, Date: %s, Début: %s, Fin: %s\n",
			reservation.ID, reservation.RoomID, reservation.Date, reservation.StartTime, reservation.EndTime)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erreur lors de l'itération sur les réservations : %v", err)
	}

	// Offre des options de navigation après avoir visualisé les réservations.
	menulogic.NavigationOptions(db, scanner)
}

// Fonction pour récupérer et afficher les réservations par date
func ViewReservationsByDate(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Print("Entrez la date pour laquelle vous souhaitez voir les réservations (format YYYY-MM-DD) : ")
	scanner.Scan()
	date := scanner.Text()

	// Validation de la date (simple vérification de format)
	if _, err := time.Parse("2006-01-02", date); err != nil {
		fmt.Println("Erreur : Format de date invalide.")
		return
	}

	// Appel à la fonction pour obtenir les réservations
	reservations, err := GetReservationsByDate(db, date)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des réservations :", err)
		return
	}

	// Affichage des réservations
	if len(reservations) == 0 {
		fmt.Println("Aucune réservation trouvée pour la date", date)
		return
	}

	fmt.Println("Réservations pour la date", date)
	for _, reservation := range reservations {
		fmt.Printf("ID Réservation: %d, ID Salle: %d, Début: %s, Fin: %s\n", reservation.ID, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	}
}

func ReservationExists(db *sql.DB, reservationID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM reservations WHERE id = ?)", reservationID).Scan(&exists)
	if err != nil {
		log.Printf("Erreur lors de la vérification de l'existence de la réservation: %v", err)
		return false
	}
	return exists
}

func GetAllReservations(db *sql.DB) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := "SELECT id, room_id, date, start_time, end_time FROM reservations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		rowErr := rows.Close()
		if rowErr != nil {
			log.Printf("Erreur: %v", rowErr)
		}
	}(rows)

	for rows.Next() {
		var r models.Reservation
		if err := rows.Scan(&r.ID, &r.RoomID, &r.Date, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	return reservations, nil
}
