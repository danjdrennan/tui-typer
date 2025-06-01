package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"strings"
	"time"
	"typr/game"
)

type TUITest struct {
	app       *tview.Application
	gameState *game.GameState
	textView  *tview.TextView
	statsView *tview.TextView
	typedText string
	allText   string
}

func NewTUITest() *TUITest {
	return &TUITest{
		app:       tview.NewApplication(),
		typedText: "",
	}
}

func (t *TUITest) RunTypingTest(gameState *game.GameState) {
	t.gameState = gameState

	// Join all words into one text string
	t.allText = strings.Join(gameState.Words, " ")

	// Create the main flex container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Stats view at top
	t.statsView = tview.NewTextView()
	t.statsView.SetBorder(true).SetTitle(" Typing Test Stats ")
	t.statsView.SetDynamicColors(true)
	t.statsView.SetBorderPadding(0, 0, 1, 1)

	// Text view in center
	t.textView = tview.NewTextView()
	t.textView.SetBorder(true).SetTitle(" Type this text ")
	t.textView.SetDynamicColors(true).SetWordWrap(true)
	t.textView.SetBorderPadding(1, 1, 2, 2)

	// Instructions at bottom
	instructions := tview.NewTextView()
	instructions.SetBorder(true).SetTitle(" Instructions ")
	instructions.SetText("Type the text above. No backspace corrections allowed. Press ESC to exit.")
	instructions.SetTextAlign(tview.AlignCenter)
	instructions.SetBorderPadding(0, 0, 1, 1)

	// Layout: stats (3 rows), text (flexible), instructions (3 rows)
	flex.AddItem(t.statsView, 5, 0, false).
		AddItem(t.textView, 0, 1, false).
		AddItem(instructions, 5, 0, false)

	// Set up input capture
	t.app.SetInputCapture(t.handleInput)

	// Start the game
	gameState.Start()

	// Start timer for updates
	go t.updateLoop()

	// Update display initially
	t.updateDisplay()

	// Run the TUI
	if err := t.app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func (t *TUITest) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC, tcell.KeyCtrlD:
		// Force exit
		t.app.Stop()
		os.Exit(0)
		return nil

	case tcell.KeyEscape:
		t.app.Stop()
		return nil

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		// Ignore backspace events - no corrections allowed
		return nil

	case tcell.KeyRune:
		char := event.Rune()

		// Handle quit commands even during test
		if char == 'q' && t.gameState.Finished {
			t.app.Stop()
			return nil
		}

		if t.gameState.Finished {
			return nil // Ignore other input after test completion
		}

		if char >= 32 && char <= 126 { // Printable characters
			t.typedText += string(char)
			t.gameState.ProcessChar(char)
			t.updateDisplay()
		}
		return nil
	}

	return event
}

func (t *TUITest) updateLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for !t.gameState.Finished {
		select {
		case <-ticker.C:
			if t.gameState.IsTimeUp() {
				t.gameState.Finish()
				t.app.QueueUpdateDraw(func() {
					t.updateDisplay()
				})
				return
			}
			t.app.QueueUpdateDraw(func() {
				t.updateDisplay()
			})
		}
	}
}

func (t *TUITest) updateDisplay() {
	// Update stats
	timeLeft := t.gameState.TestDuration - t.gameState.GetElapsedTime()
	if timeLeft < 0 {
		timeLeft = 0
	}

	statsText := fmt.Sprintf(
		"[#f9e2af]Time: [#cdd6f4]%.1fs   [#f9e2af]WPM: [#cdd6f4]%.1f   [#f9e2af]Accuracy: [#cdd6f4]%.1f%%   [#f9e2af]Errors: [#cdd6f4]%d   [#f9e2af]Progress: [#cdd6f4]%.1f%%",
		timeLeft.Seconds(),
		t.gameState.CalculateWPM(),
		t.gameState.CalculateAccuracy(),
		t.gameState.Errors,
		t.gameState.GetProgress(),
	)

	if t.gameState.Finished {
		statsText += "\n[#f38ba8]TEST COMPLETED! Press ESC/q to exit, Ctrl+C to force quit."
	}

	t.statsView.SetText(statsText)

	// Update text with overlay
	t.updateTextOverlay()
}

func (t *TUITest) updateTextOverlay() {
	var result strings.Builder
	lineWidth := 90
	currentLineLength := 0

	for i, char := range t.allText {
		// Add line breaks at word boundaries when approaching width limit
		if currentLineLength >= lineWidth && char == ' ' {
			result.WriteString("\n")
			currentLineLength = 0
			continue
		}

		if i < len(t.typedText) {
			typedChar := rune(t.typedText[i])
			if typedChar == char {
				// Correct character - Catppuccin green text
				result.WriteString(fmt.Sprintf("[#a6e3a1]%c[-]", char))
			} else {
				// Incorrect character - Catppuccin red text
				result.WriteString(fmt.Sprintf("[#f38ba8]%c[-]", char))
			}
		} else if i == len(t.typedText) {
			// Current cursor position - block character background
			result.WriteString(fmt.Sprintf("[#181825:#cdd6f4]%c[#cdd6f4:-]", char))
		} else {
			// Untyped text - Catppuccin muted
			result.WriteString(fmt.Sprintf("[#6c7086]%c[-]", char))
		}

		currentLineLength++
		if char == '\n' {
			currentLineLength = 0
		}
	}

	// Handle extra characters typed beyond the text
	if len(t.typedText) > len(t.allText) {
		extraText := t.typedText[len(t.allText):]
		for _, char := range extraText {
			result.WriteString(fmt.Sprintf("[#f38ba8]%c[-]", char))
		}
	}

	t.textView.SetText(result.String())
}

