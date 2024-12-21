package model

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type LogRegression struct {
	w [][]float64
	b []float64
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
	weights := data

	return weights
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func predict(weights [][]float64, input []float64) float64 {
	z := 0.0
	for i := 0; i < len(weights[0]); i++ {
		z += weights[0][i] * input[i] // Dot product
	}
	z += weights[1][0] // Add bias term
	return sigmoid(z)
}

func TestModel() {
	weights := LoadModel("model_weights.json")
	fmt.Println("weights length: ", len(weights.w))
	fmt.Println("bias length: ", len(weights.b))

	// TODO: load real MNIST data to test
	// input := make([]float64, 784)
	// prediction := predict(weights, input)
	// fmt.Printf("Prediction: %.2f\n", prediction)
}
