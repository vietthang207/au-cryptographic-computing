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
	B float64   `json:"b"`
}

func LoadModel() LogRegression {
	wd, _ := os.Getwd()
	weightsFile, _ := os.ReadFile(path.Join(wd, "/ppml/model/model_weights.json"))
	var data LogRegression
	// fmt.Println(json.Valid(weightsFile))
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
	z += weights.B // Add bias term
	// fmt.Println(z)
	return sigmoid(z)
}

func printWeightRange(weights LogRegression) {
	minWeight := 0.0
	maxWeight := 0.0
	for i := 0; i < len(weights.W); i++ {
		minWeight = math.Min(minWeight, weights.W[i])
		maxWeight = math.Max(maxWeight, weights.W[i])
	}

	minWeight = math.Min(minWeight, weights.B)
	maxWeight = math.Max(maxWeight, weights.B)
	fmt.Println("Min weight: ", minWeight)
	fmt.Println("Max weight: ", maxWeight)
}

func TestModel() {
	weights := LoadModel()
	fmt.Println("weights length: ", len(weights.W))
	fmt.Println("bias: ", weights.B)
	printWeightRange(weights)

	images, labels := datasets.LoadTestsetFor0And1()
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
