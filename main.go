package main

import (
	"bufio"
	"fmt"
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

		scanner.Scan()
		choice := scanner.Text()

		// Traitement de l'option choisie
		switch choice {
		case "1":
			fmt.Println("Liste des salles disponibles...")
			// Implémenter
		case "2":
			fmt.Println("Création d'une réservation...")
			// Implémenter
		case "3":
			fmt.Println("Annulation d'une réservation...")
			// Implémenter
		case "4":
			fmt.Println("Visualisation des réservations...")
			// Implémenter
		case "5":
			fmt.Println("Merci d'avoir utilisé le service. À bientôt !")
			return
		default:
			fmt.Println("Option non valide. Veuillez choisir une option entre 1 et 5.")
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}
