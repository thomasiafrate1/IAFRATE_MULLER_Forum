package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Utilisateur struct {
	IDUser        int
	Username      string
	Mail          string
	Mdp           string
	Sexe          string
	Nom           string
	Prenom        string
	DateNaissance mysql.NullTime
	DateCreation  mysql.NullTime
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var mesDonnees []Utilisateur
	for rows.Next() {
		var u Utilisateur
		if err := rows.Scan(&u.IDUser, &u.Username, &u.Mail, &u.Mdp, &u.Sexe, &u.Nom, &u.Prenom, &u.DateNaissance, &u.DateCreation); err != nil {
			log.Fatal(err)
		}
		mesDonnees = append(mesDonnees, u)
	}

	// Afficher les données
	for _, d := range mesDonnees {
		var dateNaissance string
		if d.DateNaissance.Valid {
			dateNaissance = d.DateNaissance.Time.Format("2006-01-02")
		} else {
			dateNaissance = "NULL"
		}

		var dateCreation string
		if d.DateCreation.Valid {
			dateCreation = d.DateCreation.Time.Format("2006-01-02")
		} else {
			dateCreation = "NULL"
		}

		fmt.Println(d.IDUser, d.Username, d.Mail, d.Mdp, d.Sexe, d.Nom, d.Prenom, dateNaissance, dateCreation)
	}
}
