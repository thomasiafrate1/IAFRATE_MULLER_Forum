package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password string
}

var (
	db          *sql.DB
	sessionName = "forumSession"
	store       = sessions.NewCookieStore([]byte("super-secret-key"))
)

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/forum_muller_iafrate")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	css := http.FileServer(http.Dir("../../FRONTEND/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", css))

	http.HandleFunc("/forum", forumHandler)
	http.HandleFunc("/accueil.html", forumAccueil)
	http.HandleFunc("/discussion.html", forumDiscussion)

	// ...

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	// Vérifier si l'utilisateur est connecté
	userID := session.Values["userID"]
	if userID != nil {
		// L'utilisateur est connecté, rediriger vers la page d'accueil
		http.Redirect(w, r, "/acceuil.html", http.StatusFound)
		return
	}

	// L'utilisateur n'est pas connecté, afficher le formulaire de connexion
	tmpl := template.Must(template.ParseFiles("../../FRONTEND/html/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := getUserByUsername(username)
		if err != nil {
			log.Println(err)
			http.Error(w, "Erreur lors de l'authentification", http.StatusInternalServerError)
			return
		}

		if user != nil && checkPassword(user.Password, password) {
			session, _ := store.Get(r, sessionName)
			session.Values["userID"] = user.ID
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	delete(session.Values, "userID")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func getUserByID(id int) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id_users, username, password FROM users WHERE id_users=?", id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getUserByUsername(username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id_users, username, password FROM users WHERE username=?", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Aucun utilisateur trouvé avec ce nom d'utilisateur
		}
		return nil, err
	}

	return &user, nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
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
