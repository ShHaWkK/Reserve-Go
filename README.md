# Projet de service de réservation en ligne en GO
_Alexandre UZAN
Merwane DIFALLAH
2i2_

[![MySQL](https://img.shields.io/badge/mysql-4479A1.svg?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com/fr/)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)



## Fonctionnalités
### _CLI_

- Gestion d'une base de données impliquant des réservations dans des salles en ligne de commandes
- Operations CRUD sur les salles et les réservations
- Lister les salles qui sont disponibles (la disponibilité peut être filtrée en fonction d'une date et horaires donnés si spécifié)
- Visualisation des réservation
- Récupérer les réservations par salle et par date
- Génération d'exports CSV et JSON

### _Web_

- Gestion d'une base de données impliquant des réservations dans des salles via une interface web
- Même fonctionnalités que sur la version CLI
- Génération automatique d'exports au lancement
- Possibilité de générer ces exports après avoir effectué des opérations via l'interface web puis de les télécharger

## Instructions
### _CLI_

Après avoir utilisé ``docker compose up`` et lancé le programme via ``go run main.go``, le programme se lance et affiche un menu en lignes de commandes. 
L'utilisateur a alors la possibilité de choisir une option en entrant dans le terminal le chiffre correspondant à l'option du menu que l'utilisateur souhaite exécuter.
L'utilisateur doit ensuite se laisser guider pour naviguer via le menu et a la possibilité d'entrer des champs de texte pour intéragir avec la base de données selon les options sélectionnées.

### _WEB_

Après avoir utilisé ``docker compose up`` et lancé le programme via ``go run main.go``, le programme se lance. Pour accéder à l'interface web, l'utilisateur doit se connecter au port 8095 (par défaut) du localhost. 
L'utilisateur arrive sur une page d'accueil avec des boutons correspondant aux actions qu'il est possible de faire.
En cliquant sur les boutons, l'utilisateur est amené à saisir dans des champs de texte pour les opérations qui impliquent une action de l'utilisateur pour modifier la base de données. Pour des opérations de consultation, l'utilisateur peut être amené à spécifier une salle ou une date pour filtrer les données.

## Répartition des tâches 

| Alexandre        | Merwane     
| ------|-----
| Création des fonctions majeures	| 	Structuration du code
| Gestion de la logique docker|Création des packages
| html/css de l'app web  	| 	Ajout de fonctions mineures
| Mise en place de la BDD 	| 	Rédaction

## Documentation technique

### Prérequis
- Installation des dépndances MySQL ``go get -u github.com/go-sql-driver/mysql``
### Structure du programme et Fonctionnalités
1. Définition des packages :
    - ``bufio``, ``csv``, ``json``.... : Manipulation des fichiers ainsi que des formats de données
    - ``database/sql``, ``github.com.go-sql-driver/mysql`` : Gestion de la base de données
    - 	Packages locaux :
        ``"Reserve-Go/dtb"`` : Contient le code relatif à la connexion à la BDD.
	    ``"Reserve-Go/exportlogic"`` : Contient la logique nécessaire à l'exportation des données de la BDD sous format json ou csv
	    ``"Reserve-Go/reservationlogic"`` : Contient les fonctions relatives à la manipulation des réservations
        ``"Reserve-Go/roomlogic"`` : Contient les fonctions relatives à la manipulation et opérations CRUD sur les sales
	    ``"Reserve-Go/utils"`` : Contient les fonctions pour colorer le texte et effacer l'écran pour la version CLI et les fonctions qui gèrent la redirection vers les pages de la version web.
2. Définition des structures
    - ``Room`` : Cette structure contient des informations sur les salles (ID, Name, Capacity)
    - ``Reservation`` : Cette structure contient des informations sur les réservations (ID, RoomID, Date, StartTime, EndTime)
 3. Connexion à la base de données :

    - Le programme initialise une connection à la base de données mySQL
    - La fonction ConnectToDB() est utilisée pour établir cette connexion.
    
4. Menu interactif et interface web:

   - Le programme propose un menu interactif en ligne de commande où l'utilisateur peut choisir différentes actions (par exemple: lister les salles, ajouter une salle, créer ou annuler une réservation)
   - La version web du programme propose l'équivalent du menu interfactif CLI via une interface web sur laquelle on accéde aux actions en cliquant sur des boutons

 5. Gestion des réservations :
 
    - Fonctions pour lister, ajouter, et modifier les salles, ainsi que pour créer, visualiser, et annuler les réservations
    - Fonctions pour visualiser les réservations par salle ou par date

6. Exportation de données :

    - Fonctions pour exporter les réservations en format CSV ou JSON, enregistrant les fichiers localement.
    - Pour la version WEB, la fonction permettant de télécharger les exports va d'abord appeler une fonction de génération avant de télécharger les exports au format correspondant
