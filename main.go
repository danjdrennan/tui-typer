package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"typr/data"
	"typr/game"
	"typr/ui"
)

func main() {
	wordBank, err := data.LoadWords("words.txt")
	if err != nil {
		log.Fatalf("Error loading words: %v", err)
	}

	setupSignalHandling()

	for {
		choice := ui.ShowMainMenu()

		switch choice {
		case ui.StartTest:
			runTypingTest(wordBank)
		case ui.ViewStats:
			showStats()
		case ui.Exit:
			fmt.Println("Thanks for using Typr!")
			return
		}
	}
}

func setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nExiting Typr...")
		os.Exit(0)
	}()
}


func runTypingTest(wordBank *data.WordBank) {
	for {
		words := wordBank.GenerateSequence(300)
		duration := 60 * time.Second

		gameState := game.NewGame(words, duration)
		tuiTest := ui.NewTUITest()

		tuiTest.RunTypingTest(gameState)
		showTestResults(gameState)

		// Show post-test menu
		if !ui.ShowPostTestMenu(wordBank) {
			break // Return to main menu
		}
		// Continue loop for another test
	}
}


func showTestResults(gameState *game.GameState) {
	result := gameState.GetTestResult()

	fmt.Println("\n=== Test Results ===")
	fmt.Printf("WPM: %.2f\n", result.WPM)
	fmt.Printf("Accuracy: %.2f%%\n", result.Accuracy)
	fmt.Printf("Time: %v\n", result.TestDuration)
	fmt.Printf("Words completed: %d\n", result.TotalWords)
	fmt.Printf("Errors: %d\n", result.Errors)

	err := data.SaveTestResult(result)
	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
	} else {
		fmt.Println("Results saved!")
	}
}

func showStats() {
	ui.ShowStatsMenu()
}
