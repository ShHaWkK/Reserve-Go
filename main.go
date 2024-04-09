package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	var db *sql.DB
	var err error

	for i := 0; i < 5; i++ { // Réessayer 5 fois
		connectionString := fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?parseTime=true",
			"root", // ou vos variables d'environnement
			"rootpassword",
			"projetgo",
		)

		db, err = sql.Open("mysql", connectionString)
		if err != nil {
			log.Printf("Failed to open database connection: %v", err)
			time.Sleep(2 * time.Second) // Attendre 2 secondes avant de réessayer
			continue
		}

		err = db.Ping()
		if err == nil {
			break
		}

		log.Printf("Failed to ping database: %v", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("error verifying connection to the database: %v", err)
	}

	log.Println("Successfully connected to the database.")
	return db, nil
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
	fmt.Println("Liste des salles disponibles...")
	// Implémentez la logique pour lister les salles disponibles
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
	// Écrivez la requête pour vérifier la disponibilité de la salle
	query := `SELECT COUNT(*) FROM reservations WHERE room_id = ? AND date = ? AND 
              NOT (start_time >= ? OR end_time <= ?)`

	var count int
	err := db.QueryRow(query, roomID, date, endTime, startTime).Scan(&count)
	if err != nil {
		log.Fatal("Erreur lors de la vérification de la disponibilité : ", err)
	}

	return count == 0
}

func insertReservation(db *sql.DB, roomID, date, startTime, endTime string) {
	query := `INSERT INTO reservations (room_id, date, start_time, end_time) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, roomID, date, startTime, endTime)
	if err != nil {
		log.Fatal("Erreur lors de la création de la réservation : ", err)
	}

	fmt.Println("Réservation créée avec succès.")
}

func cancelReservation(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Annulation d'une réservation...")
	// Implémentez la logique pour annuler une réservation
}

func viewReservations(db *sql.DB) {
	fmt.Println("Visualisation des réservations...")
	// Implémentez la logique pour afficher les réservations
}
