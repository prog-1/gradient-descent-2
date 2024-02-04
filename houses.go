package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type House struct {
	Square    float64
	Type      string
	Price     float64
	WallColor string
}

func readHousesFromCSV(path string) ([]House, error) {
	file, err1 := os.Open(path)
	if err1 != nil {
		log.Fatalf("Can't open file with path: %v", path)
	}
	defer file.Close()
	houses := []House{}

	reader := csv.NewReader(file)
	reader.Comma = ','

	for i := 0; ; {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if i == 0 {
			i++
			continue
		}

		square, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, err
		}

		houseType := record[1]

		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		wallColor := record[3]

		house := House{
			Square:    square,
			Type:      houseType,
			Price:     price,
			WallColor: wallColor,
		}

		houses = append(houses, house)
	}

	return houses, nil
}
