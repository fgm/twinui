package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/gorilla/mux"
	"github.com/rivo/tview"

	"github.com/fgm/twinui/model"
	"github.com/fgm/twinui/tview"
	"github.com/fgm/twinui/web"
)

func initModel(path string) (*model.Story, error) {
	story := make(model.Story)
	err := story.Load(path)
	if err != nil {
		return nil, fmt.Errorf(`loading story: %w`, err)
	}
	return &story, nil
}

// logger is a decorator for a standard logger, ensuring the text UI is refreshed
// on log writes.
type logger struct {
	app *tview.Application
	io.Writer
}

// Write implements io.Writer, by delegating writes to the underlying logger,
// but ensuring the UI gets refreshed with the written log.
func (l logger) Write(p []byte) (n int, err error) {
	ch := make(chan interface{}, 2)
	l.app.QueueUpdateDraw(func() {
		n, err = l.Writer.Write(p)
		ch <- n
		ch <- err
	})
	n = (<-ch).(int)
	err, _ = (<-ch).(error)
	return
}

func initTextUI(story *model.Story) *tview.Application {
	view := ui.NewView(story)
	app := tview.NewApplication()
	app.SetRoot(view.Grid, true).
		EnableMouse(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEsc:
				app.Stop()
			case tcell.KeyRune:
				u := view.URLFromKey(event)
				switch {
				case u == `quit`:
					app.Stop()
				case u != ``:
					view.Handle(u)
				}
			}
			return nil
		})

	view.Handle(`intro`)
	log.SetOutput(&logger{app: app, Writer: view.Body})
	return app
}

func initWebUI(story *model.Story) *mux.Router {
	tpl, err := template.ParseFiles("./web/arc.gohtml")
	if err != nil {
		log.Println("Failed parsing arc template", err)
		return nil
	}

	r := mux.NewRouter()
	r.HandleFunc("/style.css", web.StyleHandler)
	r.HandleFunc("/arc/{arc}", func(w http.ResponseWriter, r *http.Request) {
		web.ArcHandler(w, r, story, tpl)
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/arc/intro", http.StatusMovedPermanently)
	})
	return r
}

func main() {
	path := flag.String(`story`, `./model/gopher.json`, `The name of the JSON file containing the story data`)
	port := flag.Int("port", 8080, "The TCP port on which to listen")
	flag.Parse()

	// Initialize model: without data, we can't proceed.
	story, err := initModel(*path)
	if err != nil {
		log.Fatalf("Starting model: %v\n", err)
	}
	defer story.Close()

	// Initialize the twin UIs.
	app := initTextUI(story)
	router := initWebUI(story)

	// Run the twin UIs, exiting the app whenever either of them exits.
	done := make(chan bool)
	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("Running text app: %v\n", err)
		}
		done <- true
	}()

	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(*port), router); err != nil {
			log.Fatalf("Running web app: %v\n", err)
		}
		done <- true
	}()

	<-done
}
