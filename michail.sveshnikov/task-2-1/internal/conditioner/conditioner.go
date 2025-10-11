package conditioner

import (
	"errors"
)

var ErrConditioner = errors.New("invalid comparison sign")

type Conditioner struct {
	LowBound int
	UpBound  int
}

func (cond *Conditioner) OptimizeTemperature(compSign string, temperature int) (int, error) {
	const invalidTemperature int = -1

	switch compSign {
	case ">=":
		if temperature <= cond.UpBound {
			if temperature > cond.LowBound {
				cond.LowBound = temperature
			}

			return cond.LowBound, nil
		} else {
			cond.LowBound = temperature
			return invalidTemperature, nil
		}
	case "<=":
		if temperature >= cond.LowBound {
			if temperature < cond.UpBound {
				cond.UpBound = temperature
			}

			return cond.LowBound, nil
		} else {
			cond.UpBound = temperature
			return invalidTemperature, nil
		}
	default:
		return invalidTemperature, ErrConditioner
	}
}
