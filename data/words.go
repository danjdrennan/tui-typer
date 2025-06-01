package data

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Word struct {
	Text      string
	Frequency int
}

type WordBank struct {
	Words []Word
	Total int
}

func LoadWords(filename string) (*WordBank, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open words file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	wb := &WordBank{
		Words: make([]Word, 0, len(records)),
		Total: 0,
	}

	for _, record := range records {
		if len(record) != 2 {
			continue
		}

		text := record[0]
		freq, err := strconv.Atoi(record[1])
		if err != nil {
			continue
		}

		wb.Words = append(wb.Words, Word{
			Text:      text,
			Frequency: freq,
		})
		wb.Total += freq
	}

	return wb, nil
}

func (wb *WordBank) SelectRandomWord() string {
	if len(wb.Words) == 0 {
		return ""
	}

	target := rand.Intn(wb.Total)
	current := 0

	for _, word := range wb.Words {
		current += word.Frequency
		if current > target {
			return word.Text
		}
	}

	return wb.Words[len(wb.Words)-1].Text
}

func (wb *WordBank) GenerateSequence(count int) []string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	sequence := make([]string, count)

	for i := 0; i < count; i++ {
		sequence[i] = wb.SelectRandomWord()
	}

	return sequence
}

