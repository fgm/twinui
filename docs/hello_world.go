package main

import (
	"github.com/rivo/tview"
)

func main() {
	tv := tview.NewButton("Hello, world!")
	tview.NewApplication().SetRoot(tv, true).Run()
}
