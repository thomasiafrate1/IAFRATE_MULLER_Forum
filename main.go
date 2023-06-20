package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var sessionCookieName = "user_session"

type User struct {
	Username   string
	Password   string
	email      string
	sexe       string
	name       string
	first_name string
	birth_date string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/muller-iafrate-forum")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/login", loginFormHandler)
	http.HandleFunc("/home", homeFormHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/discussion", createDiscussionHandler)
	http.HandleFunc("/createurs", creatorHandler)
	http.HandleFunc("/saveMessage", saveMessageHandler)

	http.ListenAndServe(":9000", nil)

}

func handler(mux *http.ServeMux) {
	mux.HandleFunc("/css/main.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "main.css")
	})

	mux.HandleFunc("/css/connexion.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "connexion.css")
	})

	mux.HandleFunc("/connexion.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "connexion.js")
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations du formulaire
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	sexe := r.FormValue("sexe")
	name := r.FormValue("name")
	first_name := r.FormValue("first_name")
	birth_date := r.FormValue("birth_date")

	// Vérifier si l'utilisateur existe déjà dans la base de données
	if userExists(username) {
		http.Error(w, "Nom d'utilisateur déjà utilisé", http.StatusBadRequest)
		return
	}

	// Insérer l'utilisateur dans la base de données
	err := insertUser(username, email, password, sexe, name, first_name, birth_date)
	if err != nil {
		log.Println("Erreur lors de l'enregistrement:", err)
		http.Error(w, "Erreur lors de l'enregistrement", http.StatusInternalServerError)
		return
	}

	// L'utilisateur est enregistré avec succès
	// Vous pouvez effectuer d'autres actions ici, par exemple, définir une session ou rediriger vers une page d'accueil

	http.Redirect(w, r, "/login", http.StatusFound)
}

func insertUser(username, mail, password, sexe, name, first_name, birth_date string) error {
	_, err := db.Exec("INSERT INTO users (username, mail, password, sexe, name, first_name, birth_date) VALUES (?, ?, ?, ?, ?, ?, ?)", username, mail, password, sexe, name, first_name, birth_date)
	return err
}

func loginRegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.URL.Path == "/login" {
			loginHandler(w, r)
		} else if r.URL.Path == "/register" {
			registerHandler(w, r)
		}
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations du formulaire
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/muller-iafrate-forum")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Exécuter une requête pour vérifier les informations de connexion
	row := db.QueryRow("SELECT id_users FROM users WHERE username=? AND password=?", username, password)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		log.Println("Échec de la connexion:", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// L'utilisateur est connecté avec succès
	// Définir un cookie de session avec l'`id_users`
	http.SetCookie(w, &http.Cookie{
		Name:  sessionCookieName,
		Value: strconv.Itoa(userID), // convert userID to string
	})

	http.Redirect(w, r, "/home", http.StatusFound)
}

func loginFormHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifiez la méthode de la requête
	if r.Method != http.MethodPost {
		// Affichez le formulaire de connexion
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
		return
	}

	// Le formulaire a été soumis, appelez la fonction de gestion de la soumission du formulaire
	loginHandler(w, r)
}

func creatorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl1 := template.Must(template.ParseFiles("créateurs.html"))
	tmpl1.Execute(w, nil)
}

func userExists(username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", username).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'utilisateur:", err)
		return false
	}

	return count > 0
}

func authenticateUser(username, password string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=? AND password=?", username, password).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de l'authentification:", err)
		return false
	}

	return count > 0
}

/* CATEGORIE */

type Category struct {
	Id    string
	Genre string
}

func getCategories() ([]Category, error) {
	rows, err := db.Query("SELECT id_cat, gender FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.Id, &cat.Genre); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
func createDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifiez la méthode de la requête
	if r.Method != http.MethodPost {
		// Si ce n'est pas une requête POST, renvoyez une erreur 405 (Méthode non autorisée)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Récupérez le cookie de la session
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Error(w, "Vous devez être connecté pour créer une discussion", http.StatusUnauthorized)
		return
	}

	// Récupérez les valeurs du formulaire
	nameDiscussion := r.FormValue("name_discussion")
	dateStart := r.FormValue("date_start")

	// Insérez les données dans la base de données
	_, err = db.Exec("INSERT INTO discussion (name_discussion, date_start, id_users) VALUES (?, ?, ?)", nameDiscussion, dateStart, cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl1 := template.Must(template.ParseFiles("discussion.html"))
	tmpl1.Execute(w, nil)
	// Si tout va bien, redirigez vers la page discussion.html
	http.Redirect(w, r, "/discussion", http.StatusSeeOther)
}

func saveMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifiez la méthode de la requête
	if r.Method != http.MethodPost {
		// Si ce n'est pas une requête POST, renvoyez une erreur 405 (Méthode non autorisée)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Récupérez le cookie de la session
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Error(w, "Vous devez être connecté pour poster un message", http.StatusUnauthorized)
		return
	}

	log.Println("ID utilisateur récupéré à partir du cookie de session :", cookie.Value)

	// Récupérez la valeur du formulaire
	content := r.FormValue("messageInput")

	// Insérez les données dans la base de données
	result, err := db.Exec("INSERT INTO message (contained, id_users) VALUES (?, ?)", content, cookie.Value)
	if err != nil {
		log.Println("Erreur lors de l'insertion du message dans la base de données :", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Résultat de l'insertion du message :", result)
}

func someHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le cookie de session
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		// Pas de cookie de session, l'utilisateur n'est pas connecté
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	// Récupérer l'`id_users` du cookie de session
	idUsers, err := strconv.Atoi(sessionCookie.Value)
	if err != nil {
		// Valeur de cookie invalide, traiter l'erreur
		http.Error(w, "Session invalide", http.StatusInternalServerError)
		return
	}

	// À ce stade, vous avez l'ID de l'utilisateur à partir du cookie de session.
	// Vous pouvez l'utiliser pour récupérer d'autres informations sur l'utilisateur dans la base de données.

	// Exécuter une requête pour obtenir les informations de l'utilisateur
	row := db.QueryRow("SELECT username, email FROM users WHERE id_users=?", idUsers)
	var username, email string
	err = row.Scan(&username, &email)
	if err != nil {
		log.Println("Échec de la récupération des informations de l'utilisateur :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	// Créez un struct pour détenir les informations de l'utilisateur
	userInfo := struct {
		Username string
		Email    string
	}{
		Username: username,
		Email:    email,
	}

	// Afficher les informations de l'utilisateur (par exemple, dans une page HTML)
	tmpl := template.Must(template.ParseFiles("userinfo.html"))
	tmpl.Execute(w, userInfo)
}

type Discussion struct {
	Id   string
	Name string
}

func getDiscussions() ([]Discussion, error) {
	rows, err := db.Query("SELECT id_discussion, name_discussion FROM discussion")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discussions []Discussion
	for rows.Next() {
		var d Discussion
		if err := rows.Scan(&d.Id, &d.Name); err != nil {
			return nil, err
		}
		discussions = append(discussions, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return discussions, nil
}

func homeFormHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := getCategories()
	if err != nil {
		http.Error(w, "Impossible de charger les catégories", http.StatusInternalServerError)
		return
	}

	discussions, err := getDiscussions()
	if err != nil {
		http.Error(w, "Impossible de charger les discussions", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories  []Category
		Discussions []Discussion
	}{
		Categories:  categories,
		Discussions: discussions,
	}

	tmpl := template.Must(template.ParseFiles("acceuil.html"))
	tmpl.Execute(w, data)
}
