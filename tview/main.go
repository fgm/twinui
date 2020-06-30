// Demo code for the Grid primitive.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func initModel(path string) (*Story, error) {
	story := make(Story)
	err := story.Load(path)
	if err != nil {
		return nil, fmt.Errorf(`loading story: %w`, err)
	}
	return &story, nil
}

func initUI(story *Story) (*tview.Application, *View) {
	view := NewView(story)
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

	return app, view
}

func main() {
	path := flag.String(`story`, `./gophercises_cyoa/gopher.json`, `The name of the JSON file containing the story data`)
	flag.Parse()

	story, err := initModel(*path)
	if err != nil {
		log.Fatalf("Starting model: %v\n", err)
	}

	app, view := initUI(story)
	view.Handle(`intro`)
	if err := app.Run(); err != nil {
		log.Fatalf("Running app: %v\n", err)
	}
}
