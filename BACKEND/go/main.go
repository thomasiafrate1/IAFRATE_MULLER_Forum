// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"

// 	_ "github.com/go-sql-driver/mysql"
// )

// func main() {
// 	// ...

// 	r := mux.NewRouter()
// 	r.HandleFunc("/data", sendDataHandler).Methods("GET")
// 	http.Handle("/", r)

// 	// ...
// }

// func main1() {
// 	// ...
// 	rows, err := db.Query("SELECT * FROM ma_table")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var id int
// 		var name string
// 		err := rows.Scan(&id, &name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(id, name)
// 	}
// 	if err = rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// var db *sql.DB

// type Data struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// func sendDataHandler(w http.ResponseWriter, r *http.Request) {
// 	// ...

// 	rows, err := db.Query("SELECT * FROM ma_table")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	var data []Data

// 	for rows.Next() {
// 		var d Data
// 		err := rows.Scan(&d.ID, &d.Name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		data = append(data, d)
// 	}

// 	if err = rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

//		w.Header().Set("Content-Type", "application/json")
//		w.Write(jsonData)
//	}
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Utilisateur struct {
	IDUser        int
	Username      string
	Mail          string
	Mdp           string
	Sexe          string
	Nom           string
	Prenom        string
	DateNaissance time.Time // supposant que la date de naissance est stockée comme une date
	DateCreation  time.Time // supposant que la date de création est stockée comme une date
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

	// Juste pour tester l'affichage des données
	for _, d := range mesDonnees {
		fmt.Println(&d.IDUser, &d.Username, &d.Mail, &d.Mdp, &d.Sexe, &d.Nom, &d.Prenom, &d.DateNaissance, &d.DateCreation)
	}
}
