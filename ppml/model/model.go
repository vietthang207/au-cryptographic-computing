package model

import (
	"PPML/ppml/datasets"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
)

type LogRegression struct {
	W []float64 `json:"w"`
	B []float64 `json:"b"`
}

func LoadModel(weightsFilename string) LogRegression {
	weightsFile, _ := os.ReadFile(weightsFilename)
	var data LogRegression
	fmt.Println(json.Valid(weightsFile))
	err := json.Unmarshal(weightsFile, &data)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return data
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func predict(weights LogRegression, input []float64) float64 {
	z := 0.0
	for i := 0; i < len(weights.W); i++ {
		z += weights.W[i] * input[i] // Dot product
	}
	z += weights.B[0] // Add bias term
	return sigmoid(z)
}

func TestModel() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	weights := LoadModel(path.Join(wd, "/ppml/model/model_weights.json"))
	fmt.Println("weights length: ", len(weights.W))
	fmt.Println("bias length: ", len(weights.B))

	// TODO: load real MNIST data to test
	images, labels := datasets.LoadTestsetFor01()
	numTest := len(images)
	correctPrediction := 0
	for i := 0; i < numTest; i++ {
		prediction := int(math.Round(predict(weights, images[i])))
		if prediction == labels[i] {
			correctPrediction += 1
		}
	}
	accuracy := float64(correctPrediction) / float64(numTest)
	fmt.Println("Accuracy: ", accuracy)
}
