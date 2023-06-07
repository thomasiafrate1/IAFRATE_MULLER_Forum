// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	_ "github.com/go-sql-driver/mysql"
// )

// type User struct {
// 	Username string
// 	Password string
// }

// func main() {
// 	http.HandleFunc("/", loginHandler)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		// Récupérer les données du formulaire
// 		username := r.FormValue("username")
// 		password := r.FormValue("password")

// 		// Vérifier les informations d'identification dans la base de données
// 		user, err := authenticateUser(username, password)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, "Erreur lors de la vérification des informations d'identification", http.StatusInternalServerError)
// 			return
// 		}

// 		if user != nil {
// 			// Connexion réussie
// 			fmt.Fprintf(w, "Connexion réussie pour l'utilisateur: %s", username)
// 			return
// 		}

// 		// Échec de la connexion
// 		fmt.Fprint(w, "Échec de la connexion")
// 		return
// 	}

// 	// Afficher le formulaire de connexion
// 	fmt.Fprint(w, getLoginFormHTML())
// }

// func authenticateUser(username, password string) (*User, error) {
// 	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	// Requête pour vérifier l'utilisateur et le mot de passe dans la base de données
// 	query := "SELECT username, password FROM users WHERE username = ? AND password = ?"
// 	row := db.QueryRow(query, username, password)

// 	var user User
// 	err = row.Scan(&user.Username, &user.Password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// Aucun utilisateur correspondant trouvé
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
//}

func getLoginFormHTML() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>Connexion</title>
</head>
<body>
    <form action="/" method="POST">
        <label for="username">Nom d'utilisateur:</label>
        <input type="text" id="username" name="username" required>

        <label for="password">Mot de passe:</label>
        <input type="password" id="password" name="password" required>

        <button type="submit">Connexion</button>
    </form>
</body>
</html>
`
}
