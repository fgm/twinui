package web

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/fgm/twinui/model"
)

// ArcHandler handles path /arc/{arc}
func ArcHandler(w http.ResponseWriter, r *http.Request, story *model.Story, tpl *template.Template) {
	name, ok := mux.Vars(r)["arc"]
	if !ok {
		log.Println("arcHandler called without an arc")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	arc := story.Arc(name)
	if arc == nil {
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

// StyleHandler handles path /style.css
func StyleHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("web/style.css")
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
