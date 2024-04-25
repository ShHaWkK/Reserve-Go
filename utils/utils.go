package utils

import "fmt"

const (
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
)

func ColorLog(color, message string) {
	fmt.Println(color, message, ColorReset)
}

// Fonction pour colorer le texte
func ColorString(color, message string) string {
	return color + message + ColorReset
}

// Fonction pour effacer l'écran
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
