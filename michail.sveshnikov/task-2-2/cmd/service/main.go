package main

import (
	"fmt"
)

func main() {
	var numDishes int

	_, err := fmt.Scanln(&numDishes)
	if err != nil || numDishes < 1 || numDishes > 10000 {
		fmt.Println("ERROR: invalid number of dishes!")

		return
	}

	for range numDishes {
		var preferenceLevel int

		_, err = fmt.Scan(&preferenceLevel)
		if err != nil || preferenceLevel < -10000 || preferenceLevel > 10000 {
			fmt.Println("ERROR: invalid preference level of the dish!")

			return
		}
	}

	var num int

	_, err = fmt.Scanln(&num)
	if err != nil || num < 1 || num > numDishes {
		fmt.Println("ERROR: invalid number of the kth dish by preference!")

		return
	}
}
