package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"typr/game"
	"unsafe"
)

const MaxInputBuffer = 1600

type TUI struct {
	gameState     *game.GameState
	displayLines  int
	typedText     string
	totalTypedPos int
}

func NewTUI() *TUI {
	return &TUI{
		displayLines:  0,
		typedText:     "",
		totalTypedPos: 0,
	}
}

func (t *TUI) RunTypingTest(gameState *game.GameState) {
	t.gameState = gameState

	fmt.Println("\n=== Typing Test ===")
	fmt.Printf("Type the following text. Test duration: %v\n", gameState.TestDuration)
	fmt.Println("Press Enter when ready to start...")

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	t.enableRawMode()
	defer t.disableRawMode()

	gameState.Start()
	t.clearScreen()
	t.displayTestInterface()

	go t.handleTimer()
	t.handleInput()

	t.disableRawMode()
	t.clearScreen()
	fmt.Printf("\nTest completed!\n")
}

func (t *TUI) handleTimer() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for !t.gameState.Finished {
		select {
		case <-ticker.C:
			if t.gameState.IsTimeUp() {
				t.gameState.Finish()
				return
			}
			t.updateDisplay()
		}
	}
}

func (t *TUI) handleInput() {
	for !t.gameState.Finished {
		char := t.readChar()
		if char == 0 {
			continue
		}

		if char == 27 {
			fmt.Print("\033[2K\r")
			fmt.Print("Test aborted. Press any key to continue...")
			t.readChar()
			return
		}

		if char == 127 || char == 8 {
			t.handleBackspace()
		} else if char == 32 {
			t.handleSpace()
		} else if char >= 32 && char <= 126 {
			t.handleCharInput(char)
		}

		t.updateDisplay()
	}
}

func (t *TUI) handleCharInput(char rune) {
	if len(t.typedText) < MaxInputBuffer-1 {
		t.typedText += string(char)
		t.totalTypedPos++
		t.gameState.ProcessChar(char)
	}
}

func (t *TUI) handleSpace() {
	if len(t.typedText) < MaxInputBuffer-1 {
		t.typedText += " "
		t.totalTypedPos++
		t.gameState.ProcessChar(' ')
	}
}

func (t *TUI) handleBackspace() {
	if len(t.typedText) > 0 {
		t.typedText = t.typedText[:len(t.typedText)-1]
		t.totalTypedPos--
		if t.totalTypedPos < 0 {
			t.totalTypedPos = 0
		}
		t.gameState.ProcessChar('\b')
	}
}


func (t *TUI) displayTestInterface() {
	t.moveCursor(1, 1)
	fmt.Print("\033[2K")

	timeLeft := t.gameState.TestDuration - t.gameState.GetElapsedTime()
	if timeLeft < 0 {
		timeLeft = 0
	}

	fmt.Printf("Time: %.1fs | Progress: %.1f%% | WPM: %.1f | Accuracy: %.1f%% | Errors: %d",
		timeLeft.Seconds(),
		t.gameState.GetProgress(),
		t.gameState.CalculateWPM(),
		t.gameState.CalculateAccuracy(),
		t.gameState.Errors)

	t.moveCursor(3, 1)
	t.displayTextWithOverlay()

	t.displayLines = 15
}

func (t *TUI) displayTextWithOverlay() {
	allText := ""
	for i, word := range t.gameState.Words {
		if i > 0 {
			allText += " "
		}
		allText += word
	}

	for i := 3; i <= 10; i++ {
		t.moveCursor(i, 1)
		fmt.Print("\033[2K")
	}

	t.moveCursor(3, 1)
	charPos := 0
	linePos := 3
	colPos := 1

	for i, char := range allText {
		if colPos > 80 {
			linePos++
			colPos = 1
			t.moveCursor(linePos, colPos)
		}

		if i < len(t.typedText) {
			typedChar := rune(t.typedText[i])
			if typedChar == char {
				fmt.Printf("\033[42m%c\033[0m", char)
			} else {
				fmt.Printf("\033[41m%c\033[0m", char)
			}
		} else if i == len(t.typedText) {
			fmt.Printf("\033[43m%c\033[0m", char)
		} else {
			fmt.Printf("\033[90m%c\033[0m", char)
		}

		charPos++
		colPos++

		if linePos > 8 {
			break
		}
	}

	if len(t.typedText) > len(allText) {
		extraText := t.typedText[len(allText):]
		for _, char := range extraText {
			if colPos > 80 {
				linePos++
				colPos = 1
				t.moveCursor(linePos, colPos)
			}
			fmt.Printf("\033[41m%c\033[0m", char)
			colPos++
		}
	}

	t.moveCursor(10, 1)
	fmt.Printf("Typed: %d chars", len(t.typedText))
}

func (t *TUI) updateDisplay() {
	t.clearLines()
	t.displayTestInterface()
}

func (t *TUI) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func (t *TUI) clearLines() {
	for i := 0; i < t.displayLines; i++ {
		t.moveCursor(i+1, 1)
		fmt.Print("\033[2K")
	}
}

func (t *TUI) moveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}


func (t *TUI) enableRawMode() {
	exec.Command("stty", "raw", "-echo").Run()
	fmt.Print("\033[?25l")
}

func (t *TUI) disableRawMode() {
	exec.Command("stty", "cooked", "echo").Run()
	fmt.Print("\033[?25h")
}

func (t *TUI) readChar() rune {
	var buf [1]byte
	n, err := syscall.Read(int(os.Stdin.Fd()), (*(*[1]byte)(unsafe.Pointer(&buf[0])))[:])
	if err != nil || n == 0 {
		return 0
	}
	return rune(buf[0])
}
