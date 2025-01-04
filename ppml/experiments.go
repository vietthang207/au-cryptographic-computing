package ppml

import (
	"PPML/ppml/datasets"
	"PPML/ppml/model"
	"fmt"
	"math"
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

func Main() {
	// x := FloatToBigDec(-7, 100)
	// y := FloatToBigDec(3, 100)
	// z := Mul(x, y)
	// fmt.Println("x: ", x.integral)
	// fmt.Println("y: ", y.integral)
	// fmt.Println("x.y: ", z.integral)
	// fmt.Println("x.y float: ", z.ToFloat())
	model.TestModel()
	modelWeights := model.LoadModel()

	mnistTestImages, mnistTestLabels := datasets.LoadTestsetFor0And1()
	numTestImages := len(mnistTestImages)

	mnistCircuit := PolyToDotProductCircuit(MNIST_IMG_SIZE + 1) //+1 due to the dummy pixel using to multiply with the bias
	d := initDealer(mnistCircuit)

	wrongCounter := 0
	for i := 0; i < numTestImages; i++ {
		expected := mnistTestLabels[i]
		actual := simulateDotProductProtocol(mnistCircuit, mnistTestImages[i], modelWeights, d)
		if actual != expected {
			wrongCounter += 1
		}
	}
	accuracy := 1 - float64(wrongCounter)/float64(numTestImages)
	fmt.Println("Accuracy of dot product MPC: ", accuracy)
}
