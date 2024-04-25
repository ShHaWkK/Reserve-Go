package menulogic

import (
	"Reserve-Go/utils"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// ----------------------- Affichage de Menu -----------------------//
func ShowMenu() {
	//clearScreen()
	fmt.Println(utils.ColorString(utils.ColorBlue, strings.Repeat("-", 50)))
	fmt.Println(utils.ColorString(utils.ColorGreen, "Bienvenue dans le Service de Réservation en Ligne"))
	fmt.Println(utils.ColorString(utils.ColorBlue, strings.Repeat("-", 50)))
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
	fmt.Println("12. Lister les salles disponibles à un temps donné")
	fmt.Println("13. Quitter")
	fmt.Print("\nChoisissez une option : ")
}

// ------------------------------	Menu Aide 	------------------------------//
func ShowHelp() {
	utils.ClearScreen()
	fmt.Println(utils.ColorString(utils.ColorGreen, "Aide :"))
	fmt.Println(utils.ColorString(utils.ColorBlue, strings.Repeat("-", 25)))
	fmt.Println("1. Lister les salles - Affiche toutes les salles disponibles.")
	fmt.Println("2. Modifier une salle - Nous pouvons modifier les salles existantes.")
	fmt.Println("3. Créer une Salle - - Il faut entrer les informations nécessaires.")
	fmt.Println("4. Créer une réservation - Il faut entrer les informations nécessaires.")
	fmt.Println("5. Annuler une réservation - Vous aurez besoin de l'ID de la réservation.")
	fmt.Println("6. Visualiser les réservations - Pour voir les réservations existantes.")
	fmt.Println("7. Récupérer les réservation par salle ")
	fmt.Println("8. Récupérer les réservations par date")
	fmt.Println("9. Aide -> C'est nous YOUPI ! ")
	fmt.Println("10. Exportation CSV ")
	fmt.Println("11. Exportation JSON ")
	fmt.Println("12. Lister les salles disponibles à un temps donné - Entrer une date et affiche les salles disponibles à ce moment")
	fmt.Println("13. Quitter - Pour fermer l'application.")
	fmt.Println("\nAppuyez sur 'Entrée' pour retourner au menu principal.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func NavigationOptions(db *sql.DB, scanner *bufio.Scanner) {
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
