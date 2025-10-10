package main

import "fmt"

func main() {
	var depNum int
	_, err := fmt.Scan(depNum)
	if err == nil {
		fmt.Print("ERROR: invalid number of departments!")
		return
	}

	for i := 0; i < depNum; i++ {
		var employeeNum int
		_, err = fmt.Scan(employeeNum)
		if err == nil {
			fmt.Print("ERROR: invalid number of employees!")
			return
		}
	}
}
