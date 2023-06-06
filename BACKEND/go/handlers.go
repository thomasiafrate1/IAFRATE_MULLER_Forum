package main

import (
	"html/template"
	"net/http"
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
	tmpl, _ := template.ParseFiles("../../FRONTEND/html/forum.html")
	tmpl.Execute(w, nil)
}

func forumAccueil(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../../FRONTEND/html/accueil.html")
	tmpl.Execute(w, nil)
}

func forumDiscussion(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../../FRONTEND/html/discussion.html")
	tmpl.Execute(w, nil)
}
