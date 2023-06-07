package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	css := http.FileServer(http.Dir("../../FRONTEND/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", css))

	http.HandleFunc("/", forumAccueil)
	http.HandleFunc("/connexion", forumHandler)
	http.HandleFunc("/discussion", forumDiscussion)

	http.ListenAndServe(":9000", nil)
}

func forumHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../../FRONTEND/html/index.html")
	tmpl.Execute(w, nil)
}

func forumAccueil(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../../FRONTEND/html/accueil.html")
	tmpl.Execute(w, nil)
}

func forumDiscussion(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../FRONTEND/html/discussion.html")
	tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erreur de traitement du template", http.StatusInternalServerError)
		return
	}
}

type User struct {
	Username string
	Password string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Récupérer les données du formulaire
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Vérifier les informations d'identification dans la base de données
		user, err := authenticateUser(username, password)
		if err != nil {
			log.Println(err)
			http.Error(w, "Erreur lors de la vérification des informations d'identification", http.StatusInternalServerError)
			return
		}

		if user != nil {
			// Connexion réussie
			fmt.Fprintf(w, "Connexion réussie pour l'utilisateur: %s", username)
			return
		}

		// Échec de la connexion
		fmt.Fprint(w, "Échec de la connexion")
		return
	}

	// Afficher le formulaire de connexion
	fmt.Fprint(w, getLoginFormHTML())
}

func authenticateUser(username, password string) (*User, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Requête pour vérifier l'utilisateur et le mot de passe dans la base de données
	query := "SELECT username, password FROM users WHERE username = ? AND password = ?"
	row := db.QueryRow(query, username, password)

	var user User
	err = row.Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// Aucun utilisateur correspondant trouvé
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
func getLoginFormHTML() string {
	return `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
		<link rel="stylesheet" href="../assets/css/forum.css">
		<link rel="stylesheet" href="../assets/css/connexion.css">
	</head>
	<body>
	
		<video id="background-video" autoplay loop muted>
			<source src="../assets/video/traphub2- Trim.mp4" type="video/mp4">
		</video>
	
		<div id="wsh">
			<div class="vatefaireenculernemo">
				<h1 class="titre">TrapHub</h1>
				<p class="titredessous">Vivez la musique, créez des souvenirs.</p>
				<div class="boutons">
					<button id="boutonSeConnecter">Se connecter</button>
					<button id="boutonSinscrire">S'inscrire</button>
				</div>
			</div>
		</div>
	
		<div id="popupSeConnecter">
			<button id="boutonFermer">X</button>
			<h1 class="titreConnexion">TrapHub</h1>
			<div class="lesInputs">
				<form action="/login" method="POST">
					<label for="username">Nom d'utilisateur:</label>
					<input type="text" id="username" name="username" required>
			
					<label for="password">Mot de passe:</label>
					<input type="password" id="password" name="password" required>
			
					<button type="submit">Connexion</button>
				</form>
			</div>
		</div>
	
		<div id="popupSinscrire">
			<button id="boutonFermer2">X</button>
			<h1 class="titreConnexion">TrapHub</h1>
			<div class="lesInputsInscription">
				<input type="text" class="inputInscriptionUsername" placeholder="Username" autocomplete="off" >
				<input type="text" class="inputInscriptionMDP" placeholder="Password" autocomplete="off" >
				<input type="email" class="inputInscriptionEmail" placeholder="E-mail" autocomplete="off" >
			</div>
			<button class="boutonInscription">S'inscrire</button>
		</div>
	
	
	</body>
		<script src="../assets/js/connexion.js"></script>
	</html>
	`

}
