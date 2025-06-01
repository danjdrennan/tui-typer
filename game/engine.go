package game

import (
	"time"
)

type GameStatus int

const (
	StatusMenu GameStatus = iota
	StatusTyping
	StatusFinished
	StatusPaused
)

type GameState struct {
	Words          []string
	CurrentWordIdx int
	CurrentCharIdx int
	UserInput      string
	StartTime      time.Time
	TestDuration   time.Duration
	Errors         int
	TotalChars     int
	CorrectChars   int
	Status         GameStatus
	Finished       bool
}

type TestResult struct {
	Timestamp    time.Time
	WPM          float64
	Accuracy     float64
	TestDuration time.Duration
	TotalWords   int
	Errors       int
	TotalChars   int
}

func NewGame(words []string, duration time.Duration) *GameState {
	return &GameState{
		Words:          words,
		CurrentWordIdx: 0,
		CurrentCharIdx: 0,
		UserInput:      "",
		TestDuration:   duration,
		Errors:         0,
		TotalChars:     0,
		CorrectChars:   0,
		Status:         StatusMenu,
		Finished:       false,
	}
}

func (g *GameState) Start() {
	g.StartTime = time.Now()
	g.Status = StatusTyping
}

func (g *GameState) ProcessChar(char rune) {
	if g.Status != StatusTyping {
		return
	}

	if g.IsTimeUp() {
		g.Finish()
		return
	}

	currentWord := g.GetCurrentWord()
	if currentWord == "" {
		return
	}

	switch char {
	case ' ':
		g.processSpace()
	default:
		g.processTypedChar(char)
	}
}

func (g *GameState) processSpace() {
	// Space character counting
	g.TotalChars++
	g.CorrectChars++ // Space is always correct if we reach this point
	
	g.nextWord()
}


func (g *GameState) processTypedChar(char rune) {
	currentWord := g.GetCurrentWord()
	
	if g.CurrentCharIdx < len(currentWord) {
		expectedChar := rune(currentWord[g.CurrentCharIdx])
		
		g.UserInput += string(char)
		g.CurrentCharIdx++
		g.TotalChars++
		
		if char == expectedChar {
			g.CorrectChars++
		} else {
			g.Errors++
		}
	} else {
		g.UserInput += string(char)
		g.TotalChars++
		g.Errors++
	}
}

func (g *GameState) nextWord() {
	g.CurrentWordIdx++
	g.CurrentCharIdx = 0
	g.UserInput = ""
	
	if g.CurrentWordIdx >= len(g.Words) {
		g.Finish()
	}
}

func (g *GameState) GetCurrentWord() string {
	if g.CurrentWordIdx >= len(g.Words) {
		return ""
	}
	return g.Words[g.CurrentWordIdx]
}

func (g *GameState) IsTimeUp() bool {
	return time.Since(g.StartTime) >= g.TestDuration
}

func (g *GameState) Finish() {
	g.Status = StatusFinished
	g.Finished = true
}

func (g *GameState) GetElapsedTime() time.Duration {
	if g.StartTime.IsZero() {
		return 0
	}
	return time.Since(g.StartTime)
}

func (g *GameState) CalculateWPM() float64 {
	elapsed := g.GetElapsedTime()
	if elapsed == 0 {
		return 0
	}
	
	minutes := elapsed.Minutes()
	if minutes == 0 {
		return 0
	}
	
	return float64(g.CorrectChars/5) / minutes
}

func (g *GameState) CalculateAccuracy() float64 {
	if g.TotalChars == 0 {
		return 100.0
	}
	
	return float64(g.CorrectChars) / float64(g.TotalChars) * 100.0
}

func (g *GameState) GetTestResult() TestResult {
	return TestResult{
		Timestamp:    time.Now(),
		WPM:          g.CalculateWPM(),
		Accuracy:     g.CalculateAccuracy(),
		TestDuration: g.GetElapsedTime(),
		TotalWords:   g.CurrentWordIdx,
		Errors:       g.Errors,
		TotalChars:   g.TotalChars,
	}
}

func (g *GameState) GetProgress() float64 {
	if len(g.Words) == 0 {
		return 0
	}
	return float64(g.CurrentWordIdx) / float64(len(g.Words)) * 100.0
}