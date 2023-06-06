package main

import (
	"html/template"
	"net/http"
)

func main() {

	css := http.FileServer(http.Dir("../../FRONTEND/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", css))

	http.HandleFunc("/", forumHandler)
	http.ListenAndServe(":9000", nil)
}

func forumHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../../FRONTEND/html/accueil.html")
	tmpl.Execute(w, nil)
}
