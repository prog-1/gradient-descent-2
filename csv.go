package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type house struct {
	square    float64
	houseType string
	price     float64
	wallColor string
}

func readHouses(file io.Reader) (houses []house) {
	reader := csv.NewReader(file) //initializing new reader
	reader.Comma = ','

	_, _ = reader.Read() //Ommiting first line

	for {

		//Reading one record at a time
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		//Saving square as float64
		square, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		//Saving houseType as string
		houseType := record[1]

		//Saving price as float64
		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		//Saving wallColor as string
		wallColor := record[3]

		//Creating house instance
		house := house{square, houseType, price, wallColor}

		//Adding house instance to others
		houses = append(houses, house)
	}
	return houses
}

func getHouses() []house {

	//Opening csv file
	file, err := os.Open("house_prices.csv")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	//Declaring file closure at the end
	defer file.Close()

	//Saving houses data
	houses := readHouses(file)
	return houses
}
