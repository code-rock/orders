package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func home_page(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fsefffe")
	tmpl, _ := template.ParseFiles("templates/home_page.html")
	tmpl.Execute(w, "'dffgfdg'")
}

func handleRequest() {
	http.HandleFunc("/", home_page)
	http.ListenAndServe(":3001", nil)
}

func Listen() {
	handleRequest()
}
