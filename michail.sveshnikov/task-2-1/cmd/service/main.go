package main

import (
	"fmt"

	"github.com/svesh3000/task-2-1/internal/conditioner"
)

func main() {
	const (
		lowLimit int = 15
		upLimit  int = 30
	)

	var depNum int

	_, err := fmt.Scan(&depNum)
	if err != nil || depNum < 1 || depNum > 1000 {
		fmt.Println("ERROR: invalid number of departments!")

		return
	}

	for range depNum {
		var employeeNum int

		_, err = fmt.Scan(&employeeNum)
		if err != nil || employeeNum < 1 || employeeNum > 1000 {
			fmt.Println("ERROR: invalid number of employees!")

			continue
		}

		cond := conditioner.Conditioner{
			LowBound: lowLimit,
			UpBound:  upLimit,
		}

		for range employeeNum {
			var compSign string

			_, err = fmt.Scan(&compSign)
			if err != nil {
				fmt.Println("ERROR: invalid comparison sign!")

				continue
			}

			var temperature int

			_, err = fmt.Scan(&temperature)
			if err != nil {
				fmt.Println("ERROR: invalid temperature!")

				continue
			}

			optimizeTemp, err := cond.OptimizeTemperature(compSign, temperature)
			if err != nil {
				fmt.Println(err)

				continue
			}
			fmt.Println(optimizeTemp)
		}
	}
}
