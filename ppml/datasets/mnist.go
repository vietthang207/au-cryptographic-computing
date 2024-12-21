package datasets

import (
	"fmt"

	mnist "github.com/petar/GoMNIST"
)

func rawImageToFloatVector(rawImage mnist.RawImage) []float64 {
	res := make([]float64, len(rawImage))
	for i := 0; i < len(rawImage); i++ {
		res[i] = float64(rawImage[i]) / 256.0
	}
	return res
}

func LoadTestset() ([][]float64, []int) {
	_, test, err := mnist.Load("./data")
	if err != nil {
		panic(err)
	}
	testSetSize := test.Count()
	fmt.Printf("Load %d test images\n", testSetSize)
	images := make([][]float64, test.Count())
	for i := 0; i < testSetSize; i++ {
		images[i] = rawImageToFloatVector(test.Images[i])
	}

	labels := make([]int, testSetSize)
	for i := 0; i < testSetSize; i++ {
		labels[i] = int(test.Labels[i])
	}
	return images, labels
}

// Only using 0 and 1 digits
func LoadTestsetFor0And1() ([][]float64, []int) {
	images, labels := LoadTestset()

	images01 := make([][]float64, 0, 2000)
	labels01 := make([]int, 0, 2000)

	for i := 0; i < len(images); i++ {
		if labels[i] == 0 || labels[i] == 1 {
			images01 = append(images01, images[i])
			labels01 = append(labels01, labels[i])
		}
	}
	fmt.Printf("Load %d test images for 0 and 1\n", len(images01))
	return images01, labels01
}
