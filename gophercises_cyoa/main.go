package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	path := flag.String("story", "gopher.json", "The name of the JSON file containing the story data")
	port := flag.Int("port", 8080, "The TCP port on which to listen")
	flag.Parse()

	story, err := loadStory(*path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	tpl, err := template.ParseFiles("arc.gohtml")
	if err != nil {
		log.Println("Failed parsing arc template", err)
		os.Exit(2)
	}

	r := mux.NewRouter()
	r.HandleFunc("/style.css", styleHandler)
	r.HandleFunc("/arc/{arc}", func(w http.ResponseWriter, r *http.Request) {
		arcHandler(w, r, story, tpl)
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/arc/intro", http.StatusMovedPermanently)
	})
	_ = http.ListenAndServe(":"+strconv.Itoa(*port), r)
}
