package ppml

import (
	"PPML/ppml/datasets"
	"PPML/ppml/model"
	"fmt"
	"math"
	"time"
)

const MNIST_IMG_SIZE = 784

func simulateDotProductProtocol(circuit circuit, img []float64, modelWeights model.LogRegression, d dealer) int {
	img = append(img, 1.0) //adding dummy pixel to multiply with the bias
	a := initAlice(circuit, img, d)
	weights := append(modelWeights.W, modelWeights.B)
	b := initBob(circuit, weights, d)
	for !a.hasOutput() {
		receive(&b, send(&a))
		receive(&a, send(&b))
	}
	dotProduct := a.output()
	sigmoid := 1 / (1 + math.Exp(dotProduct))
	return int(math.Round(sigmoid))
}

func simulateTaylorSeriesProtocol(circuit circuit, img []float64, modelWeights model.LogRegression, d dealer) int {
	img = append(img, 1.0) //adding dummy pixel to multiply with the bias
	a := initAlice(circuit, img, d)
	weights := append(modelWeights.W, modelWeights.B)
	b := initBob(circuit, weights, d)
	for !a.hasOutput() {
		receive(&b, send(&a))
		receive(&a, send(&b))
	}
	sigmoid := a.output()
	return int(math.Round(sigmoid))
}

func printCircuit(circuit circuit) {
	for i := 0; i < len(circuit.gates); i++ {
		switch circuit.gates[i] {
		case InputA:
			fmt.Println(i, ": InputA")
		case InputB:
			fmt.Println(i, ": InputB")
		case Output:
			fmt.Println(i, ": Output")
		case AddConst:
			fmt.Println(i, ": AddConst")
		case MulConst:
			fmt.Println(i, ": MulConst")
		case Add2Wires:
			fmt.Println(i, ": Add2Wires")
		case Mul2Wires:
			fmt.Println(i, ": Mul2Wires")
		}
	}
}

func experimentDotProduct(modelWeights model.LogRegression, mnistTestImages [][]float64, mnistTestLabels []int) {
	numTestImages := len(mnistTestImages)

	mnistCircuit := PolyToDotProductCircuit(MNIST_IMG_SIZE + 1) //+1 due to the dummy pixel using to multiply with the bias
	d := initDealer(mnistCircuit)

	wrongCounter := 0
	start := time.Now()
	for i := 0; i < numTestImages; i++ {
		expected := mnistTestLabels[i]
		actual := simulateDotProductProtocol(mnistCircuit, mnistTestImages[i], modelWeights, d)
		if actual != expected {
			wrongCounter += 1
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %v microseconds\n", elapsed.Microseconds())
	accuracy := 1 - float64(wrongCounter)/float64(numTestImages)
	fmt.Println("Accuracy of dot product MPC: ", accuracy)

}

func experimentTaylorSeries(modelWeights model.LogRegression, mnistTestImages [][]float64, mnistTestLabels []int, approxDegree int) {
	numTestImages := len(mnistTestImages)
	mnistCircuit := PolyToCircuit(MNIST_IMG_SIZE+1, approxDegree) //+1 due to the dummy pixel using to multiply with the bias
	d := initDealer(mnistCircuit)

	// printCircuit(mnistCircuit)

	wrongCounter := 0
	start := time.Now()
	for i := 0; i < numTestImages; i++ {
		expected := mnistTestLabels[i]
		actual := simulateTaylorSeriesProtocol(mnistCircuit, mnistTestImages[i], modelWeights, d)
		if actual != expected {
			wrongCounter += 1
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %v microseconds\n", elapsed.Microseconds())
	accuracy := 1 - float64(wrongCounter)/float64(numTestImages)
	fmt.Printf("Accuracy of MPC with Taylor approximation of degree %d: %f\n", approxDegree, accuracy)
}

func Main() {
	// TestCoefficient()
	model.TestModel()
	modelWeights := model.LoadModel()

	mnistTestImages, mnistTestLabels := datasets.LoadTestsetFor0And1()

	experimentDotProduct(modelWeights, mnistTestImages, mnistTestLabels)

	for approxDegree := 1; approxDegree < 20; approxDegree += 2 {
		experimentTaylorSeries(modelWeights, mnistTestImages, mnistTestLabels, approxDegree)
	}

}
