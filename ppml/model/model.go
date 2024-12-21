package model

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

func LoadModel(weightsFilename string) [][]float64 {
	weightsFile, _ := os.ReadFile(weightsFilename)
	var weights [][]float64
	err := json.Unmarshal(weightsFile, &weights)
	if err != nil {
		panic("Unmarshal err")
	}

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
	fmt.Println("model length: ", len(weights))
	fmt.Println("weights length: ", len(weights[0]))
	fmt.Println("bias length: ", len(weights[1]))

	// TODO: load real MNIST data to test
	// input := make([]float64, 784)
	// prediction := predict(weights, input)
	// fmt.Printf("Prediction: %.2f\n", prediction)
}
