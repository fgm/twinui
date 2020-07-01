package ui

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/fgm/twinui/model"
)

// View holds the structure of the application View:
type View struct {
	// Heading is the top line.
	Heading *tview.TextView
	// Body is the main frame.
	Body    *tview.TextView
	// Actions contains the action menu.
	Actions *tview.List
	// Grid is the container wrapping the Heading, Body, and Actions.
	*tview.Grid
	// Story is the model from which the View reads data.
	*model.Story
}

// Handle updates the View from an Arc loaded from the Story by its URL.
func (v View) Handle(url string) {
	arc := v.Story.Arc(url)
	if arc ==  nil {
		log.Printf("Path not found: %s\n", url)
		return
	}
	fmt.Fprint(v.Heading.Clear(), arc.Title)
	b := v.Body.Clear()
	for _, row := range arc.Body {
		fmt.Fprintln(b, row + "\n")
	}
	v.Actions.Clear()
	if len(arc.Options) == 0 {
		arc.Options = []model.Option{{
			Label: `Leave story`,
			URL:   `quit`,
		}}
	}
	for k, item := range arc.Options {
		v.Actions.InsertItem(k, item.Label, item.URL, rune('a' + k), nil)
	}
}

// URLFromKey resolves a key event to an Arc URL if possible.
func (v View) URLFromKey(event *tcell.EventKey) string {
	if event.Key() != tcell.KeyRune {
		return ``
	}
	index := int(event.Rune() - 'a')
	if index < 0 || index >= v.Actions.GetItemCount() {
		return ``
	}
	_, url := v.Actions.GetItemText(index)
	return url
}

func textView(title string) *tview.TextView {
	tv := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorBlack)

	tv.SetBackgroundColor(tcell.ColorWhite).
		SetBorderColor(tcell.ColorLightGray).
		SetBorder(true)

	tv.SetTitle(` ` + title + ` `).
		SetTitleColor(tcell.ColorSlateGray).
		SetTitleAlign(tview.AlignLeft)
	return tv
}

func list(title string) *tview.List {
	l := tview.NewList().
		SetMainTextColor(tcell.ColorBlack).
		ShowSecondaryText(false).
		SetShortcutColor(tcell.ColorDarkGreen)

	l.SetBackgroundColor(tcell.ColorWhite).
		SetBorderColor(tcell.ColorLightGray).
		SetBorder(true).
		SetTitle(` ` + title + ` `).
		SetTitleColor(tcell.ColorSlateGray).
		SetTitleAlign(tview.AlignLeft)
	return l
}

// NewView builds an initialized full-screen View.
func NewView(story *model.Story) *View {
	v := &View{
		Heading: textView("Scene"),
		Body:    textView("Description").SetScrollable(true),
		Actions: list("Choose wisely"),
		Grid:    tview.NewGrid(),
		Story:   story,
	}
	v.Grid.
		SetRows(3, 0, 5). // 1-row title, 3-row actions. Add 2 for their own borders.
		SetBorders(false). // Use the view borders instead.
		AddItem(v.Heading, 0, 0, 1, 1, 0, 0, false).
		AddItem(v.Body, 1, 0, 1, 1, 0, 0, false).
		AddItem(v.Actions, 2, 0, 1, 1, 0, 0, true)

	return v
}
