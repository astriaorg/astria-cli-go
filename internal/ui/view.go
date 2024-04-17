package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Props is an empty interface for passing data to the view.
type Props interface{}

type View interface {
	// Render returns the tview.Flex that represents the view.
	// The Props argument is used to pass data to the view.
	Render(p Props) *tview.Flex
	// GetKeyboard is a callback for defining keyboard shortcuts
	GetKeyboard(a AppController) func(evt *tcell.EventKey) *tcell.EventKey
}

// Append the settings status to the end of the input string
func appendStatus(text string, status bool) string {
	if status {
		return text + ": [black:white]ON [-:-]"
	} else {
		return text + ": [white:darkslategray]off[-:-]"
	}
}
