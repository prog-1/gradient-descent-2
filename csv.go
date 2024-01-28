package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	houseTypeCount = 5 //count of the.. houseType's
)

type house struct {
	square    float64
	houseType string
	price     float64
	wallColor string
}

func readHouses(file io.Reader) (houses []house, houseTypes []string) {
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

		//Saving houseType in houseTypes array if it's not there
		if !contains(houseTypes, houseType) {
			houseTypes = append(houseTypes, houseType)
		}

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
	return houses, houseTypes
}

// To check whether houseType is already saved
func contains(strings []string, str string) bool {
	for i := range strings {
		if strings[i] == str {
			return true
		}
	}
	return false
}

func getHouses() ([]house, []string) {

	//Opening csv file
	file, err := os.Open("house_prices.csv")
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	//Declaring file closure at the end
	defer file.Close()

	//Saving houses data
	houses, houseTypes := readHouses(file)
	return houses, houseTypes
}

// Returns squares and prices of the houses grouped by houseType
func groupHouses() (px [houseTypeCount][]float64, py [houseTypeCount][]float64) {

	houses, houseTypes := getHouses() //getting house data from CSV file

	for _, h := range houses { //for every house
		for ht := 0; ht < houseTypeCount; ht++ { //for every houseType
			if h.houseType == houseTypes[ht] { //if type of selected house is current taken houseType
				px[ht] = append(px[ht], h.square) //saving square data to desired slice
				py[ht] = append(py[ht], h.price)  //saving price data to desired slice
			}
		}
	}
	return px, py
}
