package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "utilisateur:mdp@tcp(localhost:3306)/ma_base_de_donnees")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func main1() {
	// ...
	rows, err := db.Query("SELECT * FROM ma_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
