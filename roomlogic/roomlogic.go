package roomlogic

import (
	"Reserve-Go/menulogic"
	"Reserve-Go/models"
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

func AddRoom(db *sql.DB, scanner *bufio.Scanner) {
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
	menulogic.NavigationOptions(db, scanner)
}

func ListAvailableRooms(db *sql.DB, date string, startTime string, endTime string, scanner *bufio.Scanner) ([]models.Room, error) {
	var rooms []models.Room
	query := `SELECT id, name, capacity FROM rooms WHERE id NOT IN (
				SELECT room_id FROM reservations WHERE date = ? AND NOT (end_time <= ? OR start_time >= ?)
			) AND available = TRUE`
	rows, err := db.Query(query, date, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		rowErr := rows.Close()
		if rowErr != nil {
			log.Printf("Erreur: %v", rowErr)
		}
	}(rows)
	fmt.Printf("Salles disponnibles:")
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
		fmt.Printf("ID: %d, Nom: %s, Capacité: %d\n", room.ID, room.Name, room.Capacity)
	}
	menulogic.NavigationOptions(db, scanner)
	return rooms, nil
}

func UpdateRoom(db *sql.DB, scanner *bufio.Scanner) {
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
	menulogic.NavigationOptions(db, scanner)
}

func ListRooms(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Liste des salles disponibles:")

	query := "SELECT id, name, capacity FROM rooms"
	rows, err := db.Query(query)

	if err != nil {
		log.Printf("Erreur lors de la récupération des salles : %v", err)
		return
	}

	defer func(rows *sql.Rows) {
		rowErr := rows.Close()
		if rowErr != nil {
			log.Printf("Erreur: %v", rowErr)
		}
	}(rows)

	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Capacity); err != nil {
			log.Printf("Erreur lors de la lecture des données de la salle : %v", err)
			continue
		}
		fmt.Printf("ID: %d, Nom: %s, Capacité: %d\n", room.ID, room.Name, room.Capacity)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erreur lors de l'itération sur les salles : %v", err)
	}
	menulogic.NavigationOptions(db, scanner)
}

func IsRoomAvailable(db *sql.DB, roomID, date, startTime, endTime string) bool {
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
func IsRoomExists(db *sql.DB, roomID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = ?)", roomID).Scan(&exists)
	return err == nil && exists
}
