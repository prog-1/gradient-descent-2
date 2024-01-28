package main

import (
	"fmt"
	"time"
)

const (
	lrB    = 0.5
	lrK    = 0.000129
	epochs = 1000
)

//#######################################################################

func loss(k, b float64, px, py []float64) float64 {
	totalE := 0.0       // error sum of all points
	for i := range px { //for every point
		x, y := px[i], py[i]          //saving current point x & y for convenience
		totalE += (y - (k*x + b)) * 2 //add to total error
	}
	totalE /= float64(len(px)) //get average error
	return totalE
}

func inference(x, k, b float64) float64 {
	return k*x + b
}

func gradientDescent(k float64, b []float64, px, py [houseTypeCount][]float64, epoch int) (float64, []float64) {

	n := float64(getLenght(px[:])) //get all point count

	//########### Calculating K ##############

	dk := 0.0                                // gradient for k
	for ht := 0; ht < houseTypeCount; ht++ { //for every houseType
		for i := range px[ht] { //for every point in current houseType
			x, y := px[ht][i], py[ht][i]            //saving current point x & y for convenience
			dk -= (2 / n) * (y - (k*x + b[ht])) * x //adjusting the gradient for k
		}
	}
	k -= dk * lrK //adding gradient to k

	//########### Calculating B's ##############

	db := make([]float64, houseTypeCount)    // gradients for b coefficients for every houseType
	for ht := 0; ht < houseTypeCount; ht++ { //for every houseType
		for i := range px[ht] { //for every point in current houseType
			x, y := px[ht][i], py[ht][i]            //saving current point x & y for convenience
			db[ht] -= (2 / n) * (y - (k*x + b[ht])) //adjusting the gradient of i'th b
		}
		b[ht] -= db[ht] * lrB //adding i'th gradient to b of i'th houseType
	}

	//Printing gradients on every 100th epoch
	if epoch%100 == 0 {
		fmt.Printf("dk: %v ", dk)                //gradient for k
		for ht := 0; ht < houseTypeCount; ht++ { //for every houseType
			fmt.Printf("| db №%v: %v ", ht+1, db[ht]) //print b gradient of current houseType
		}
		fmt.Print("\n")
	}
	return k, b
}

func (a *App) linearRegression(px, py [houseTypeCount][]float64) {
	var k float64                              //initial k coefficient
	b := make([]float64, houseTypeCount)       //initial b coefficients for every houseType
	for epoch := 1; epoch <= epochs; epoch++ { //for every epoch

		//Debug print of epoch and losses
		if epoch%100 == 0 { //every 100th epoch
			fmt.Printf("\nEpoch: %v ", epoch)        //printing debug data
			for ht := 0; ht < houseTypeCount; ht++ { //for every housetype
				fmt.Printf("| Loss №%v: %v ", ht+1, loss(k, b[ht], px[ht], py[ht]))
			}
			fmt.Print("\n")
		}

		k, b := gradientDescent(k, b, px, py, epoch) //get trained values
		a.updatePlot(k, b, px, py)                   //recreating plot with new values
		time.Sleep(time.Millisecond)                 //delay to monitor the updates
	}
}

//#######################################################################

// To get all point count
func getLenght(s [][]float64) (l int) {
	for i := 0; i < len(s); i++ { //for every type
		l += len(s[i]) //sum the point count
	}
	return l
}
