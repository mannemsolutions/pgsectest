package internal

import "fmt"

type TestScore struct {
	Min    float32 `yaml:"min"`
	Max    float32 `yaml:"max"`
	Weight float32 `yaml:"weight"`
}

func (ts TestScore) FromResult(result float64) float64 {
	if ts.Max == ts.Min {
		return 0
	}
	return float64(ts.Weight) * (float64(result) - float64(ts.Min)) / (float64(ts.Max - ts.Min))
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
