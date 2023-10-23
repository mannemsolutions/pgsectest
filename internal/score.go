package internal

import "fmt"

type TestScore struct {
	Min    float32 `yaml:"min"`
	Max    float32 `yaml:"max"`
	Weight float32 `yaml:"weight"`
}

func (ts TestScore) FromResult(dividend float64, divisor float64) float64 {
	if ts.Max == ts.Min {
		return 0
	}
	if divisor == 0 {
		return float64(ts.Max)
	}
	result := ((dividend / divisor) - float64(ts.Min)) / (float64(ts.Max - ts.Min))
	if result < 0 {
		return 0
	} else if result > 1 {
		return float64(ts.Weight)
	}
	return float64(ts.Weight) * result
}

func (ts TestScore) Flawless() float64 {
	return float64(ts.Weight)
}

func (ts TestScore) Validate() (err error) {
	if ts.Min == ts.Max {
		return fmt.Errorf("Score min cannot be the same as score max")
	}
	return nil
}
