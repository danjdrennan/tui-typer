package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"typr/game"
)

const StatsFileName = "stats.csv"

func SaveTestResult(result game.TestResult) error {
	file, err := os.OpenFile(StatsFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open stats file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		result.Timestamp.Format(time.RFC3339),
		fmt.Sprintf("%.2f", result.WPM),
		fmt.Sprintf("%.2f", result.Accuracy),
		result.TestDuration.String(),
		strconv.Itoa(result.TotalWords),
		strconv.Itoa(result.Errors),
		strconv.Itoa(result.TotalChars),
	}

	return writer.Write(record)
}

func LoadTestResults() ([]game.TestResult, error) {
	file, err := os.Open(StatsFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []game.TestResult{}, nil
		}
		return nil, fmt.Errorf("failed to open stats file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	results := make([]game.TestResult, 0, len(records))

	for _, record := range records {
		if len(record) != 7 {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			continue
		}

		wpm, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}

		accuracy, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}

		duration, err := time.ParseDuration(record[3])
		if err != nil {
			continue
		}

		totalWords, err := strconv.Atoi(record[4])
		if err != nil {
			continue
		}

		errors, err := strconv.Atoi(record[5])
		if err != nil {
			continue
		}

		totalChars, err := strconv.Atoi(record[6])
		if err != nil {
			continue
		}

		results = append(results, game.TestResult{
			Timestamp:    timestamp,
			WPM:          wpm,
			Accuracy:     accuracy,
			TestDuration: duration,
			TotalWords:   totalWords,
			Errors:       errors,
			TotalChars:   totalChars,
		})
	}

	return results, nil
}
