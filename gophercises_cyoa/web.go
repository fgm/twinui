package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// arcHandler handles path /arc/{arc}
func arcHandler(w http.ResponseWriter, r *http.Request, story Story, tpl *template.Template) {
	name, ok := mux.Vars(r)["arc"]
	if !ok {
		log.Println("arcHandler called without an arc")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	arc, ok := story[name]
	if !ok {
		log.Printf("Incorrect arc requested: %s\n", name)
		http.NotFound(w, r)
		return
	}

	err := tpl.Execute(w, arc)
	if err != nil {
		log.Printf("Failed executing arc template for arc %#v: %v\n", arc, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// styleHandler handles path /style.css
func styleHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("style.css")
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/css")
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println("Failed sending CSS", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
