package main

import (
	"fmt"
	"time"
)

const (
	lrK    = 0.001
	lrB    = 0.5
	epochs = 1000
)

//#######################################################################

func gradientDescent(k []float64, b []float64, px, py [houseTypeCount][]float64, epoch int) ([]float64, []float64) {

	n := float64(getLenght(px[:])) //get all point count

	//########### Calculating coefficients ##############

	dk := make([]float64, houseTypeCount)    // gradients for k coefficients for every houseType
	db := make([]float64, houseTypeCount)    // gradients for b coefficients for every houseType
	for ht := 0; ht < houseTypeCount; ht++ { //for every houseType
		for i := range px[ht] { //for every point in current houseType
			x, y := px[ht][i], py[ht][i]                    //saving current point x & y for convenience
			dk[ht] -= (2 / n) * (y - (k[ht]*x + b[ht])) * x //adjusting the gradient of i'th k
			db[ht] -= (2 / n) * (y - (k[ht]*x + b[ht]))     //adjusting the gradient of i'th b
		}
		k[ht] -= dk[ht] * lrK //adding i'th gradient to i'th k
		b[ht] -= db[ht] * lrB //adding i'th gradient to i'th b
	}

	//########### Debug print ##############
	if epoch%100 == 0 { //every 100th epoch
		fmt.Printf("\nEpoch %v:", epoch)         //epoch number
		for ht := 0; ht < houseTypeCount; ht++ { //for every houseType
			fmt.Printf("\nHouse Type №%v ", ht+1) //houseType
			fmt.Printf("| Loss: %.4f ", loss(k[ht], b[ht], px[ht], py[ht]))
			fmt.Printf("| dk: %.4f ", dk[ht]) //gradient for k
			fmt.Printf("| db: %.4f ", db[ht]) //print b gradient of current houseType
		}
		fmt.Println()
	}
	//######################################

	return k, b
}

func (a *App) linearRegression(px, py [houseTypeCount][]float64) {
	k := make([]float64, houseTypeCount)       //initial k coefficients for every houseType
	b := make([]float64, houseTypeCount)       //initial b coefficients for every houseType
	for epoch := 1; epoch <= epochs; epoch++ { //for every epoch
		k, b := gradientDescent(k, b, px, py, epoch) //get trained values
		a.updatePlot(k, b, px, py)                   //recreating plot with new values
		time.Sleep(time.Millisecond)                 //delay to monitor the updates
	}
}

//#######################################################################

func inference(x, k, b float64) float64 {
	return k*x + b
}

func loss(k, b float64, px, py []float64) float64 {
	totalE := 0.0       // error sum of all points
	for i := range px { //for every point
		x, y := px[i], py[i]          //saving current point x & y for convenience
		totalE += (y - (k*x + b)) * 2 //add to total error
	}
	totalE /= float64(len(px)) //get average error
	return totalE
}

//#######################################################################

// To get all point count
func getLenght(s [][]float64) (l int) {
	for i := 0; i < len(s); i++ { //for every type
		l += len(s[i]) //sum the point count
	}
	return l
}
