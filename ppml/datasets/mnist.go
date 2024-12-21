package datasets

import (
	mnist "github.com/petar/GoMNIST"
)

type DataPoint struct {
	x []float64
	y float64
}

func LoadTestset() *mnist.Set {
	_, test, err := mnist.Load("./data")
	if err != nil {
		panic(err)
	}
	return test
}
