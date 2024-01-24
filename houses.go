package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

func enum(s string) []float64 {
	switch s {
	case "Duplex":
		return []float64{1, 0, 0, 0, 0}
	case "Multi-family":
		return []float64{0, 1, 0, 0, 0}
	case "Detached":
		return []float64{0, 0, 1, 0, 0}
	case "Semi-detached":
		return []float64{0, 0, 0, 1, 0}
	case "Townhouse":
		return []float64{0, 0, 0, 0, 1}
	default:
		panic("Should not happen")
	}
}

func Get() (x [][]float64, y []float64) {
	file, err := os.Open("house_prices.csv")
	if err != nil {
		log.Fatal("Can't open file: ", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		a, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			continue
		}
		b, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}
		x = append(x, append(enum(record[1]), a))
		y = append(y, b)
	}
	return
}
