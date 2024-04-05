package ui

import (
	"testing"
	"time"

	"github.com/astria/astria-cli-go/internal/testutils"
	"github.com/rivo/tview"
)

func TestProcessPane_DisplayOutput(t *testing.T) {
	mockPR := new(testutils.MockProcessRunner)

	// Setting up the mock behavior
	didStartChan := make(chan bool, 1)
	didStartChan <- true
	close(didStartChan)

	mockPR.On("GetTitle").Return("Test Process")
	mockPR.On("GetOutput").Return("Initial output").Once()
	mockPR.On("GetOutput").Return("Initial output\nUpdated output")

	app := tview.NewApplication()
	processPane := NewProcessPane(app, mockPR)
	processPane.TickerInterval = 10

	processPane.StartScan()

	go func() {
		if err := app.Run(); err != nil {
			t.Fatal("Failed to run application:", err)
		}
	}()

	// NOTE - needs to be at least double the TickerInterval to ensure the output is updated
	time.Sleep(25 * time.Millisecond)

	defer app.Stop()

	mockPR.AssertExpectations(t)

	//app.QueueUpdateDraw(func() {
	//	text := processPane.GetTextView().GetText(true)
	//	assert.Contains(t, text, "Updated output", "The textView should contain the latest process output")
	//	assert.Equal(t, 2, processPane.lineCount, "The line count should be updated")
	//})
}

//func TestProcessPane_Highlight(t *testing.T) {
//	mockPR := new(testutils.MockProcessRunner)
//
//	mockPR.On("GetTitle").Return("Test Highlight")
//
//	app := tview.NewApplication()
//	processPane := NewProcessPane(app, mockPR)
//	processPane.TickerInterval = 10
//
//	go func() {
//		if err := app.Run(); err != nil {
//			t.Fatal("Failed to run application:", err)
//		}
//	}()
//
//	// NOTE - should be a bit more than the TickerInterval
//	time.Sleep(12 * time.Millisecond) // Adjust based on need, but keep it minimal
//
//	// highlight the pane
//	app.QueueUpdateDraw(func() {
//		processPane.Highlight(true)
//	})
//
//	// NOTE - should be a bit more than the TickerInterval
//	time.Sleep(12 * time.Millisecond) // Adjust based on need
//
//	// Check if the border color is set to expected value after highlight
//	borderColor := processPane.GetTextView().GetBorderColor()
//	if borderColor != tcell.ColorBlue { // Assuming blue is the highlight color
//		t.Errorf("Expected border color to be blue when highlighted, got %v", borderColor)
//	}
//
//	// unhighlight the pane
//	app.QueueUpdateDraw(func() {
//		processPane.Highlight(false)
//	})
//
//	// NOTE - should be a bit more than the TickerInterval
//	time.Sleep(12 * time.Millisecond) // Adjust based on need
//
//	// Check if the border color is set to expected value after highlight
//	borderColor = processPane.GetTextView().GetBorderColor()
//	if borderColor != tcell.ColorGray {
//		t.Errorf("Expected border color to be blue when highlighted, got %v", borderColor)
//	}
//
//	defer app.Stop()
//}
