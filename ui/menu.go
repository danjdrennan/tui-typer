package ui

import (
	"fmt"
	"os"
	"typr/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MenuChoice int

const (
	StartTest MenuChoice = iota
	ViewStats
	Exit
)

type Menu struct {
	app      *tview.Application
	choice   MenuChoice
	selected bool
	menuView *tview.TextView
}

func NewMenu() *Menu {
	return &Menu{
		app:      tview.NewApplication(),
		choice:   StartTest,
		selected: false,
	}
}

func (m *Menu) Show() MenuChoice {
	// Create main container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// ASCII Art header
	asciiArt := tview.NewTextView()
	asciiArt.SetBorder(false)
	asciiArt.SetText(m.getASCIIArt())
	asciiArt.SetTextAlign(tview.AlignCenter)
	asciiArt.SetDynamicColors(true)

	// Menu options
	m.menuView = tview.NewTextView()
	m.menuView.SetBorder(true)
	m.menuView.SetTitle(" Menu ")
	m.menuView.SetDynamicColors(true)
	m.menuView.SetTextAlign(tview.AlignCenter)

	// Instructions
	instructions := tview.NewTextView()
	instructions.SetBorder(false)
	instructions.SetText("[#6c7086]Navigate: [#cdd6f4]↑↓ or j/k[#6c7086] | Select: [#cdd6f4]Enter/Space/1-3[#6c7086] | Exit: [#cdd6f4]ESC/q, Ctrl+C")
	instructions.SetTextAlign(tview.AlignCenter)
	instructions.SetDynamicColors(true)

	// Layout
	flex.AddItem(asciiArt, 8, 0, false).
		AddItem(m.menuView, 8, 0, false).
		AddItem(instructions, 3, 0, false)

	// Set up input handling
	m.app.SetInputCapture(m.handleMenuInput)

	// Update menu display
	m.updateMenuDisplay(m.menuView)

	// Run the menu
	if err := m.app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return m.choice
}

func (m *Menu) handleMenuInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		os.Exit(0)
	case tcell.KeyEscape:
		m.choice = Exit
		m.app.Stop()
	case tcell.KeyEnter:
		m.selected = true
		m.app.Stop()
	case tcell.KeyUp:
		m.moveUp()
	case tcell.KeyDown:
		m.moveDown()
	case tcell.KeyRune:
		switch event.Rune() {
		case 'j':
			m.moveDown()
		case 'k':
			m.moveUp()
		case 'q':
			m.choice = Exit
			m.app.Stop()
		case ' ':
			m.selected = true
			m.app.Stop()
		case '1':
			m.choice = StartTest
			m.selected = true
			m.app.Stop()
		case '2':
			m.choice = ViewStats
			m.selected = true
			m.app.Stop()
		case '3':
			m.choice = Exit
			m.selected = true
			m.app.Stop()
		}
	}
	return nil
}

func (m *Menu) moveUp() {
	if m.choice > StartTest {
		m.choice--
		m.updateDisplay()
	}
}

func (m *Menu) moveDown() {
	if m.choice < Exit {
		m.choice++
		m.updateDisplay()
	}
}

func (m *Menu) updateDisplay() {
	m.app.QueueUpdateDraw(func() {
		m.updateMenuDisplay(m.menuView)
	})
}

func (m *Menu) updateMenuDisplay(menuView *tview.TextView) {
	options := []string{
		"Start Typing Test (60 seconds)",
		"View Statistics",
		"Exit",
	}

	var menuText string
	for i, option := range options {
		if MenuChoice(i) == m.choice {
			menuText += fmt.Sprintf("[#181825:#f9e2af] > %s < [#cdd6f4:-]\n", option)
		} else {
			menuText += fmt.Sprintf("[#6c7086]   %s   [-]\n", option)
		}
	}

	menuView.SetText(menuText)
}

func (m *Menu) getASCIIArt() string {
	return `[#f9e2af]
 ████████ ██    ██ ██████  ██████
    ██     ██  ██  ██   ██ ██   ██
    ██      ████   ██████  ██████
    ██       ██    ██      ██   ██
    ██       ██    ██      ██   ██
[#6c7086]
    Terminal Typing Test
[-]`
}

func ShowMainMenu() MenuChoice {
	menu := NewMenu()
	return menu.Show()
}

func ShowPostTestMenu(wordBank *data.WordBank) bool {
	app := tview.NewApplication()

	// Create main container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Results message
	message := tview.NewTextView()
	message.SetBorder(true)
	message.SetTitle(" Test Complete! ")
	message.SetText("[#a6e3a1]Great job! Your results have been saved.\n\n[#cdd6f4]What would you like to do next?")
	message.SetTextAlign(tview.AlignCenter)
	message.SetDynamicColors(true)

	// Options
	options := tview.NewTextView()
	options.SetBorder(false)
	options.SetText("[#f9e2af]Press [#cdd6f4]SPACE[#f9e2af] for another test | Press [#cdd6f4]ESC[#f9e2af] to return to main menu")
	options.SetTextAlign(tview.AlignCenter)
	options.SetDynamicColors(true)

	// Layout
	flex.AddItem(message, 8, 0, false).
		AddItem(options, 5, 0, false)

	restart := false

	// Input handling
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			// Force exit
			app.Stop()
			os.Exit(0)
			return nil
		case tcell.KeyEscape:
			restart = false
			app.Stop()
			return nil
		case tcell.KeyRune:
			char := event.Rune()
			switch char {
			case ' ':
				restart = true
				app.Stop()
				return nil
			case 'q', 'Q':
				restart = false
				app.Stop()
				return nil
			}
		}
		return event
	})

	// Run the post-test menu
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return restart
}

func ShowStatsMenu() {
	app := tview.NewApplication()

	// Load results
	results, err := data.LoadTestResults()

	// Create main container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Header
	header := tview.NewTextView()
	header.SetBorder(true)
	header.SetTitle(" Your Statistics ")
	header.SetDynamicColors(true)
	header.SetTextAlign(tview.AlignCenter)

	var statsText string

	if err != nil {
		statsText = fmt.Sprintf("[#f38ba8]Error loading stats: %v", err)
	} else if len(results) == 0 {
		statsText = "[#f9e2af]No test results found.\n[#6c7086]Complete a typing test to see your statistics here!"
	} else {
		var totalWPM, totalAccuracy float64
		bestWPM := 0.0

		for _, result := range results {
			totalWPM += result.WPM
			totalAccuracy += result.Accuracy
			if result.WPM > bestWPM {
				bestWPM = result.WPM
			}
		}

		avgWPM := totalWPM / float64(len(results))
		avgAccuracy := totalAccuracy / float64(len(results))

		statsText = fmt.Sprintf(
			"[#a6e3a1]Total tests: [#cdd6f4]%d\n\n"+
			"[#a6e3a1]Average WPM: [#cdd6f4]%.2f\n"+
			"[#a6e3a1]Best WPM: [#cdd6f4]%.2f\n"+
			"[#a6e3a1]Average Accuracy: [#cdd6f4]%.2f%%",
			len(results), avgWPM, bestWPM, avgAccuracy)
	}

	header.SetText(statsText)

	// Recent tests
	recentView := tview.NewTextView()
	recentView.SetBorder(true)
	recentView.SetTitle(" Recent Tests ")
	recentView.SetDynamicColors(true)

	if len(results) > 0 {
		var recentText string
		start := len(results) - 8
		if start < 0 {
			start = 0
		}

		for i := start; i < len(results); i++ {
			r := results[i]
			recentText += fmt.Sprintf("[#6c7086]%s: [#cdd6f4]%.2f WPM, %.2f%% accuracy\n",
				r.Timestamp.Format("Jan 2 15:04"), r.WPM, r.Accuracy)
		}
		recentView.SetText(recentText)
	} else {
		recentView.SetText("[#6c7086]No recent tests to display")
	}

	// Instructions
	instructions := tview.NewTextView()
	instructions.SetBorder(false)
	instructions.SetText("[#6c7086]Press [#cdd6f4]ESC[#6c7086] to return to main menu")
	instructions.SetTextAlign(tview.AlignCenter)
	instructions.SetDynamicColors(true)

	// Layout
	flex.AddItem(header, 10, 0, false).
		AddItem(recentView, 0, 1, false).
		AddItem(instructions, 3, 0, false)

	// Input handling
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlD:
			// Force exit
			app.Stop()
			os.Exit(0)
			return nil
		case tcell.KeyEscape:
			app.Stop()
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' || event.Rune() == 'Q' {
				app.Stop()
				return nil
			}
		}
		return event
	})

	// Run the stats menu
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
