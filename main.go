package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// ------------------------------------------------------//
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

const (
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
)

func colorLog(color, message string) {
	fmt.Println(color, message, ColorReset)
}

//------------------------------------------------------//

func main() {
	// Connexion à la base de données
	db, err := connectToDB()
	if err != nil {
		colorLog(ColorRed, "Erreur lors de la connexion à la base de données: "+err.Error())
		return
	}
	defer db.Close()
	colorLog(ColorGreen, "Connexion à la base de données réussie.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		scanner.Scan()
		choice := scanner.Text()
		switch choice {
		case "1":
			listRooms(db)
		case "2":
			updateRoom(db, scanner)
		case "3":
			addRoom(db, scanner)
		case "4":
			createReservation(db, scanner)
		case "5":
			cancelReservation(db, scanner)
		case "6":
			viewReservations(db, scanner)
		case "7":
			viewReservationsByRoom(db, scanner)
		case "8":
			viewReservationsByDate(db, scanner)
		case "9":
			showHelp()
		case "10":
			if err := exportReservationsAsCSV(db, "reservations.csv"); err != nil {
				log.Printf("Failed to export reservations as CSV: %v", err)
			}
		case "11":
			if err := exportReservationsAsJSON(db, "reservations.json"); err != nil {
				log.Printf("Failed to export reservations as JSON: %v", err)
			}
		case "12":
			fmt.Println("Merci d'avoir utilisé le service. À bientôt !")
			return
		default:
			fmt.Println("Option non valide. Veuillez choisir une option entre 1 et 10.")
		}
	}
}

//-----------------------		Connexion 		----------------------------//

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

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			colorLog(ColorGreen, "Successfully connected to the database.")
			return db, nil
		}
		colorLog(ColorRed, "Failed to connect to the database, retrying in 1 second...")
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("error verifying connection to the database: %v", err)
}

// ----------------------- Affichage de Menu -----------------------//
func showMenu() {
	//clearScreen()
	fmt.Println(colorString(ColorBlue, strings.Repeat("-", 50)))
	fmt.Println(colorString(ColorGreen, "Bienvenue dans le Service de Réservation en Ligne"))
	fmt.Println(colorString(ColorBlue, strings.Repeat("-", 50)))
	fmt.Println("1. Lister les salles disponibles")
	fmt.Println("2. Modifier une salle")
	fmt.Println("3. Créer une Salle ")
	fmt.Println("4. Créer une réservation")
	fmt.Println("5. Annuler une réservation")
	fmt.Println("6. Visualiser les réservations")
	fmt.Println("7. Récupérer les réservations par salle")
	fmt.Println("8. Récupérer les réservations par date")
	fmt.Println("9. Aide")
	fmt.Println("10. Exportation CSV ")
	fmt.Println("11. Exportation JSON ")
	fmt.Println("12. Quitter")
	fmt.Print("\nChoisissez une option : ")
}

// ------------------------------	Menu Aide 	------------------------------//
func showHelp() {
	clearScreen()
	fmt.Println(colorString(ColorGreen, "Aide :"))
	fmt.Println(colorString(ColorBlue, strings.Repeat("-", 25)))
	fmt.Println("1. Lister les salles - Affiche toutes les salles disponibles.")
	fmt.Println("2. Créer une réservation - Il faut entrer les informations nécessaires.")
	fmt.Println("3. Annuler une réservation - Vous aurez besoin de l'ID de la réservation.")
	fmt.Println("4. Visualiser les réservations - Pour voir les réservations existantes.")
	fmt.Println("5. Quitter - Pour fermer l'application.")
	fmt.Println("\nAppuyez sur 'Entrée' pour retourner au menu principal.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Fonction pour colorer le texte
func colorString(color, message string) string {
	return color + message + ColorReset
}

// Fonction pour effacer l'écran
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

/*
*Navigation
 */
func navigationOptions(db *sql.DB, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n1. Retourner au menu principal")
		fmt.Println("2. Quitter")
		fmt.Print("\nChoisissez une option : ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			return
		case "2":
			fmt.Println("Merci d'avoir utilisé le service. À bientôt !")
			os.Exit(0)
		default:
			fmt.Println("Option non valide. Veuillez choisir une option entre 1 et 2.")
		}
	}
}

/*
* Créer  une salle
 */
func addRoom(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Ajout d'une nouvelle salle...")

	fmt.Println("Entrez le nom de la salle :")
	scanner.Scan()
	name := scanner.Text()

	fmt.Println("Entrez la capacité de la salle :")
	scanner.Scan()
	capacity, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Printf("Erreur : Capacité invalide. %v", err)
		return
	}

	query := "INSERT INTO rooms (name, capacity) VALUES (?, ?)"
	_, err = db.Exec(query, name, capacity)
	if err != nil {
		log.Printf("Erreur lors de l'ajout de la salle : %v", err)
	} else {
		fmt.Println("Salle ajoutée avec succès.")
	}
	navigationOptions(db, scanner)
}

/*
* Lister les salles disponible
 */
func listAvailableRooms(db *sql.DB, date string, startTime string, endTime string) ([]Room, error) {
	var rooms []Room

	query := `SELECT id, name, capacity FROM rooms WHERE id NOT IN (
				SELECT room_id FROM reservations WHERE date = ? AND NOT (end_time <= ? OR start_time >= ?)
			) AND available = TRUE`
	rows, err := db.Query(query, date, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

/*
* Modifier les salles
 */

func updateRoom(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Modification d'une salle existante...")

	fmt.Println("Entrez l'ID de la salle à modifier :")
	scanner.Scan()
	id, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Printf("Erreur : ID invalide. %v", err)
		return
	}

	fmt.Println("Entrez le nouveau nom (laissez vide pour ne pas modifier) :")
	scanner.Scan()
	name := scanner.Text()

	fmt.Println("Entrez la nouvelle capacité (laissez vide pour ne pas modifier) :")
	scanner.Scan()
	capacityStr := scanner.Text()
	var capacity int
	if capacityStr != "" {
		capacity, err = strconv.Atoi(capacityStr)
		if err != nil {
			log.Printf("Erreur : Capacité invalide. %v", err)
			return
		}
	}

	query := "UPDATE rooms SET name = COALESCE(NULLIF(?, ''), name), capacity = COALESCE(NULLIF(?, 0), capacity) WHERE id = ?"
	_, err = db.Exec(query, name, capacity, id)
	if err != nil {
		log.Printf("Erreur lors de la modification de la salle : %v", err)
	} else {
		fmt.Println("Salle modifiée avec succès.")
	}
	navigationOptions(db, scanner)
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
	fmt.Println(colorString(ColorBlue, strings.Repeat("-", 25)))
	fmt.Println("Création d'une réservation...")
	fmt.Println(colorString(ColorBlue, strings.Repeat("-", 25)))

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
		colorLog(ColorRed, "La salle n'est pas disponible pour le créneau demandé.")
	}
	navigationOptions(db, scanner)
}

/*
*	Récupérer les réservations par salle
 */
func viewReservationsByRoom(db *sql.DB, scanner *bufio.Scanner) {
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
	if !isRoomExists(db, roomID) {
		fmt.Println("La salle avec l'ID", roomID, "n'existe pas.")
		return
	}

	// Appel à getReservationsByRoom avec l'ID de la salle
	reservations, err := getReservationsByRoom(db, roomID)
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

/*
*	Récupérer les réservations par salle
 */
func getReservationsByRoom(db *sql.DB, roomID int) ([]Reservation, error) {
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
	return reservations, nil
}

/*
*	Récupérer les réservations par date
 */
func getReservationsByDate(db *sql.DB, date string) ([]Reservation, error) {
	var reservations []Reservation
	query := "SELECT id, room_id, date, start_time, end_time FROM reservations WHERE date = ?"
	rows, err := db.Query(query, date)
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
func isRoomExists(db *sql.DB, roomID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = ?)", roomID).Scan(&exists)
	return err == nil && exists
}

func insertReservation(db *sql.DB, roomID, date, startTime, endTime string) {
	query := `INSERT INTO reservations (room_id, date, start_time, end_time) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, roomID, date, startTime, endTime)
	if err != nil {
		log.Printf("Erreur lors de la création de la réservation : %v", err)
	} else {
		//fmt.Println("Réservation créée avec succès.")
	}

}

func cancelReservation(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Print("Entrez l'identifiant de la réservation à annuler : ")
	scanner.Scan()
	reservationID := scanner.Text()

	// Vérification de l'existence de la réservation avant de tenter de l'annuler
	if reservationExists(db, reservationID) {
		err := deleteReservation(db, reservationID)
		if err != nil {
			fmt.Println("Erreur lors de l'annulation de la réservation :", err)
		} else {
			fmt.Println("Réservation annulée avec succès.")
		}
	} else {
		fmt.Println("Aucune réservation trouvée avec cet identifiant.")
	}
	navigationOptions(db, scanner)
}

func deleteReservation(db *sql.DB, reservationID string) error {
	query := "DELETE FROM reservations WHERE id = ?"
	_, err := db.Exec(query, reservationID)
	return err
}

func viewReservations(db *sql.DB, scanner *bufio.Scanner) {
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

	// Offre des options de navigation après avoir visualisé les réservations.
	navigationOptions(db, scanner)
}

/*
*	Récupérer et afficher les réservations par date
 */

// Fonction pour récupérer et afficher les réservations par date
func viewReservationsByDate(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Print("Entrez la date pour laquelle vous souhaitez voir les réservations (format YYYY-MM-DD) : ")
	scanner.Scan()
	date := scanner.Text()

	// Validation de la date (simple vérification de format)
	if _, err := time.Parse("2006-01-02", date); err != nil {
		fmt.Println("Erreur : Format de date invalide.")
		return
	}

	// Appel à la fonction pour obtenir les réservations
	reservations, err := getReservationsByDate(db, date)
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

/*
* Vérification de l'existence d'une réservation
 */
func reservationExists(db *sql.DB, reservationID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM reservations WHERE id = ?)", reservationID).Scan(&exists)
	if err != nil {
		log.Printf("Erreur lors de la vérification de l'existence de la réservation: %v", err)
		return false
	}
	return exists
}

/*
* Exportatioon des réservations en JSON
 */
func exportReservationsAsJSON(db *sql.DB, filename string) error {
	reservations, err := getAllReservations(db)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(reservations, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

/*
*
 */
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

/*
* Exportation CSV
 */
func exportReservationsAsCSV(db *sql.DB, filename string) error {
	reservations, err := getAllReservations(db)
	if err != nil {
		log.Printf("Error fetching reservations: %v", err)
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating CSV file: %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "RoomID", "Date", "StartTime", "EndTime"}
	if err := writer.Write(header); err != nil {
		log.Printf("Error writing header to CSV: %v", err)
		return err
	}

	for _, reservation := range reservations {
		record := []string{
			strconv.Itoa(reservation.ID),
			strconv.Itoa(reservation.RoomID),
			reservation.Date,
			reservation.StartTime,
			reservation.EndTime,
		}
		if err := writer.Write(record); err != nil {
			log.Printf("Error writing record to CSV: %v", err)
			return err
		}
	}

	log.Printf("Reservations successfully exported to %s", filename)
	return nil
}
