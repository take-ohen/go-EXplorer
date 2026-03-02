package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// Run initializes and starts the GUI application.
func Run() {
	var mw *walk.MainWindow

	MainWindow{
		AssignTo: &mw,
		Title:    "Go-EXplorer",
		MinSize:  Size{Width: 800, Height: 600},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "Ready to explore...",
			},
		},
	}.Run()
}
