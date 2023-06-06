package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Message struct {
	IDMessage        int
	contient         string
	DateCreation     mysql.NullTime
	DateModification mysql.NullTime
	IDUser           int
	IDDuscussion     int
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM message")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var mesDonnees []Message
	for rows.Next() {
		var u Message
		if err := rows.Scan(&u.IDMessage, &u.contient, &u.DateModification, &u.DateCreation, &u.IDUser, &u.IDDuscussion); err != nil {
			log.Fatal(err)
		}
		mesDonnees = append(mesDonnees, u)
	}

	for _, d := range mesDonnees {
		var dateModification string
		if d.DateModification.Valid {
			dateModification = d.DateModification.Time.Format("2006-01-02")
		} else {
			dateModification = "NULL"
		}
		// Afficher les données
		var dateCreation string
		if d.DateCreation.Valid {
			dateCreation = d.DateCreation.Time.Format("2006-01-02")
		} else {
			dateCreation = "NULL"
		}

		fmt.Println(d.IDMessage, d.contient, dateModification, dateCreation, d.IDUser, d.IDDuscussion)
	}
}
