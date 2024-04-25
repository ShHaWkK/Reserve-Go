package dtb

import "database/sql"

//-------------------------- CONNEXION --------------------------//

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/projetgo")
	if err != nil {
		return nil, err
	}
	return db, nil
}
