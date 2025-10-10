package main

import "fmt"

const (
	lowLimit           int = 15
	upLimit            int = 30
	invalidTemperature int = -1
)

type Conditioner struct {
	lowBound int
	upBound  int
}

func (cond *Conditioner) optimizeTemperature(compSign string, temperature int) int {
	switch compSign {
	case ">=":
		if temperature <= cond.upBound {
			if temperature >= cond.lowBound {
				cond.lowBound = temperature
			}
			return cond.lowBound
		} else {
			return invalidTemperature
		}
	case "<=":
		if temperature >= cond.lowBound {
			if temperature <= cond.upBound {
				cond.upBound = temperature
			}
			return cond.upBound
		} else {
			return invalidTemperature
		}
	default:
		fmt.Println("ERROR: invalid comparison sign!")
	}
	return invalidTemperature
}

func main() {
	var depNum int
	_, err := fmt.Scanln(depNum)
	if err == nil || depNum < 1 || depNum > 1000 {
		fmt.Println("ERROR: invalid number of departments!")
		return
	}

	for i := 0; i < depNum; i++ {
		var employeeNum int
		_, err = fmt.Scanln(employeeNum)
		if err == nil || employeeNum < 1 || employeeNum > 1000 {
			fmt.Println("ERROR: invalid number of employees!")
			continue
		}

		cond := Conditioner{lowLimit, upLimit}
		for j := 0; j < employeeNum; j++ {
			var compSign string
			_, err = fmt.Scanln(compSign)
			if err == nil {
				fmt.Println("ERROR: invalid comparison sign!")
				return
			}
			var temperature int
			_, err = fmt.Scanln(temperature)
			if err == nil {
				fmt.Println("ERROR: invalid temperature!")
				return
			}
			fmt.Println(cond.optimizeTemperature(compSign, temperature))
		}
	}
}
