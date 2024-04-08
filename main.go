package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()

		scanner.Scan()
		choice := scanner.Text()

		// Traitement de l'option choisie
		switch choice {
		case "1":
			listRooms()
		case "2":
			createReservation()
		case "3":
			cancelReservation()
		case "4":
			viewReservations()
		case "5":
			fmt.Println("Merci d'avoir utilisé le service. À bientôt !")
			return
		default:
			fmt.Println("Option non valide. Veuillez choisir une option entre 1 et 5.")
		}
	}
}

func connectToDB() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), // Assume DB_HOST environment variable is set to your MySQL host, e.g., "localhost"
		os.Getenv("DB_PORT"), // Assume DB_PORT environment variable is set to your MySQL port, e.g., "3306"
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
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

func listRooms() {
	fmt.Println("Liste des salles disponibles...")

}

func createReservation() {
	fmt.Println("Création d'une réservation...")

}

func cancelReservation() {
	fmt.Println("Annulation d'une réservation...")
}

func viewReservations() {
	fmt.Println("Visualisation des réservations...")

}
