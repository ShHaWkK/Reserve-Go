package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"time"
)

// Structure Room
type Room struct {
	ID       int
	Name     string
	Capacity int
}

// Structure Reservation
type Reservation struct {
	ID        int
	RoomID    int
	Date      string
	StartTime string
	EndTime   string
}

func main() {
	// Connexion à la base de données
	db, err := connectToDB()
	if err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données:", err)
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			listRooms(db)
		case "2":
			createReservation(db, scanner)
		case "3":
			cancelReservation(db, scanner)
		case "4":
			viewReservations(db)
		case "5":
			fmt.Println("Merci d'avoir utilisé le service. À bientôt !")
			return
		default:
			fmt.Println("Option non valide. Veuillez choisir une option entre 1 et 5.")
		}
	}
}

func connectToDB() (*sql.DB, error) {
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

	// Attendre que la base de données soit prête
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to the database.")
			return db, nil
		}
		log.Printf("Failed to connect to the database, retrying in 1 second...")
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("error verifying connection to the database: %v", err)
}

func showMenu() {
	fmt.Print("\033[36m")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("Bienvenue dans le Service de Réservation en Ligne")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print("\033[0m")
	fmt.Println("1. Lister les salles disponibles")
	fmt.Println("2. Créer une réservation")
	fmt.Println("3. Annuler une réservation")
	fmt.Println("4. Visualiser les réservations")
	fmt.Println("5. Quitter")
	fmt.Print("\nChoisissez une option : ")
}

func listRooms(db *sql.DB) {
	fmt.Println("Liste des salles disponibles:")

	query := "SELECT id, name, capacity FROM rooms"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Erreur lors de la récupération des salles : %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
			log.Printf("Erreur lors de la lecture des données de la salle : %v", err)
			continue
		}
		fmt.Printf("ID: %d, Nom: %s, Capacité: %d\n", room.ID, room.Name, room.Capacity)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erreur lors de l'itération sur les salles : %v", err)
	}
}

func createReservation(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Création d'une réservation...")
	fmt.Println("----------------------------------")

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

	if isRoomAvailable(db, roomID, date, startTime, endTime) {
		insertReservation(db, roomID, date, startTime, endTime)
		fmt.Println("Réservation créée avec succès.")
	} else {
		fmt.Println("La salle n'est pas disponible pour le créneau demandé.")
	}
}

func isRoomAvailable(db *sql.DB, roomID, date, startTime, endTime string) bool {
	query := `SELECT COUNT(*) FROM reservations 
              WHERE room_id = ? 
                AND date = ?
                AND NOT (start_time >= ? OR end_time <= ?)`

	var count int
	err := db.QueryRow(query, roomID, date, endTime, startTime).Scan(&count)
	if err != nil {
		log.Printf("Erreur lors de la vérification de la disponibilité : %v", err)
		return false
	}

	return count == 0
}

func insertReservation(db *sql.DB, roomID, date, startTime, endTime string) {
	query := `INSERT INTO reservations (room_id, date, start_time, end_time) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, roomID, date, startTime, endTime)
	if err != nil {
		log.Printf("Erreur lors de la création de la réservation : %v", err)
	} else {
		fmt.Println("Réservation créée avec succès.")
	}
}

func cancelReservation(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Annulation d'une réservation...")
	// Implémentez la logique pour annuler une réservation
}

func viewReservations(db *sql.DB) {
	fmt.Println("Visualisation des réservations:")

	query := `SELECT r.id, r.room_id, r.date, r.start_time, r.end_time 
              FROM reservations r
              ORDER BY r.date, r.start_time`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Erreur lors de la récupération des réservations : %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var reservation Reservation
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
}
