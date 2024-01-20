package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

func Get() (x, y []float64) {
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
		x = append(x, a)
		y = append(y, b)
	}
	return
}
