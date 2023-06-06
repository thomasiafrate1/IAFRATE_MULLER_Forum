package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Discussion struct {
	IDDuscussion int
	name_dis     string
	DateCreation mysql.NullTime
	IDUser       int
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM discussion")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var mesDonnees []Discussion
	for rows.Next() {
		var u Discussion
		if err := rows.Scan(&u.IDDuscussion, &u.name_dis, &u.DateCreation, &u.IDUser); err != nil {
			log.Fatal(err)
		}
		mesDonnees = append(mesDonnees, u)
	}

	// Afficher les données
	for _, d := range mesDonnees {
		var dateCreation string
		if d.DateCreation.Valid {
			dateCreation = d.DateCreation.Time.Format("2006-01-02")
		} else {
			dateCreation = "NULL"
		}

		fmt.Println(d.IDDuscussion, d.name_dis, dateCreation, d.IDUser)
	}
}
