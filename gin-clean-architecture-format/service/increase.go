package service

import (
	"encoding/csv"
	"log"
	"os"
)

type Increase struct {
	Today     int
	Yesterday int
}

func NewIncrease(today int, yesterday int) *Increase {
	return &Increase{Today: today, Yesterday: yesterday}
}

func (increase *Increase) IncreaseRate() int {
	increaseRate := (increase.Today - increase.Yesterday) / increase.Yesterday * 100
	return increaseRate
}

func WriteCsv(records [][]string) (string, error) {
	today, err := GetToday()
	csvFileName := today + ".csv"
	if err != nil {
		panic(err)
	}
	file, err := os.Create(csvFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
			return "", err
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
		return "", err
	}
	return csvFileName, nil
}
