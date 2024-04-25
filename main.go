package main

import (
	"Reserve-Go/dtb"
	"Reserve-Go/exportlogic"
	"Reserve-Go/menulogic"
	"Reserve-Go/reservationlogic"
	"Reserve-Go/roomlogic"
	"Reserve-Go/utils"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	// Connexion à la base de données
	db, err := dtb.ConnectToDB()
	if err != nil {
		utils.ColorLog(utils.ColorRed, "Erreur lors de la connexion à la base de données: "+err.Error())
		return
	}
	defer func(db *sql.DB) {
		dbErr := db.Close()
		if dbErr != nil {
			log.Printf("Erreur: %v", dbErr)
		}
	}(db)
	utils.ColorLog(utils.ColorGreen, "Connexion à la base de données réussie.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		menulogic.ShowMenu()
		scanner.Scan()
		choice := scanner.Text()
		var date string
		var startTime string
		var endTime string
		switch choice {
		case "1":
			roomlogic.ListRooms(db, scanner)
		case "2":
			roomlogic.UpdateRoom(db, scanner)
		case "3":
			roomlogic.AddRoom(db, scanner)
		case "4":
			reservationlogic.CreateReservation(db, scanner)
		case "5":
			reservationlogic.CancelReservation(db, scanner)
		case "6":
			reservationlogic.ViewReservations(db, scanner)
		case "7":
			reservationlogic.ViewReservationsByRoom(db, scanner)
		case "8":
			reservationlogic.ViewReservationsByDate(db, scanner)
		case "9":
			menulogic.ShowHelp()
		case "10":
			if err := exportlogic.ExportReservationsAsCSV(db, "reservations.csv", scanner); err != nil {
				log.Printf("Failed to export reservations as CSV: %v", err)
			}
		case "11":
			if err := exportlogic.ExportReservationsAsJSON(db, "reservations.json", scanner); err != nil {
				log.Printf("Failed to export reservations as JSON: %v", err)
			}
		case "12":
			fmt.Println("Saisissez la date")
			fmt.Scanln(&date)
			fmt.Println("Saisissez l'heure de début")
			fmt.Scanln(&startTime)
			fmt.Println("Saisissez l'heure de fin")
			fmt.Scanln(&endTime)
			_, listErr := roomlogic.ListAvailableRooms(db, date, startTime, endTime, scanner)
			if listErr != nil {
				log.Printf("Erreur: %v", listErr)
			}
		case "13":
			fmt.Println("Merci d'avoir utilisé le service. À bientôt !")
			return
		default:
			fmt.Println("Option non valide. Veuillez choisir une option entre 1 et 10.")
		}
	}
}
