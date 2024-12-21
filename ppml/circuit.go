package ppml

type ArithGate int

const (
	InputA ArithGate = iota
	InputB
	Output
	AddConst
	MulConst
	Add2Wires
	Mul2Wires
)

type circuit struct {
	gates        []ArithGate
	firstInputs  []int
	secondInputs []int //this is index, so use int
	constants    []BigDec
}

func PolyToCircuit(inputSize int, approxDegree int) circuit {
	//Circuit encoding convention:
	//Gate                  | firstInput          | secondInput
	//Input                 | index of input bit  | 0
	//XOR/AND with constant | index of input wire | constant                   |
	//Binary gate           | index of input wire | index of second input wire |
	var gates = []ArithGate{InputA, InputA, InputA, InputB, InputB, InputB, AddConst, AddConst, AddConst, Mul2Wires, Mul2Wires, Mul2Wires, AddConst, AddConst, AddConst, Mul2Wires, Mul2Wires, Output}
	var firstInputs = []int{0, 1, 2, 0, 1, 2, 0, 1, 2, 6, 7, 8, 9, 10, 11, 12, 15, 0}
	var secondInputs = []int{0, 0, 0, 0, 0, 0, 1, 1, 1, 3, 4, 5, 1, 1, 1, 13, 14, 0}
	var constant_floats = []float64{0, 0, 0, 0, 0, 0, 1, 1, 1, 3, 4, 5, 1, 1, 1, 13, 14, 0}
	constants := FloatArrayToBigDec(constant_floats, DEFAULT_SCALAR)
	return circuit{gates: gates, firstInputs: firstInputs, secondInputs: secondInputs, constants: constants}
}

func PolyToDotProductCircuit(inputSize int) circuit {
	g1, f1, s1 := getInputSubCircuit(inputSize)
	g2, f2, s2 := getDotProductSubCircuit(inputSize)
	g := append(g1, g2...)
	f := append(f1, f2...)
	s := append(s1, s2...)
	g = append(g, Output)
	f = append(f, 0)
	s = append(s, 0)
	constants := make([]BigDec, 0)
	return circuit{gates: g, firstInputs: f, secondInputs: s, constants: constants}
}

func getInputSubCircuit(inputSize int) ([]ArithGate, []int, []int) {
	// var gates = []ArithGate{InputA, InputA, InputA, InputB, InputB, InputB}
	// var firstInputs = []int{0, 1, 2, 0, 1, 2}
	// var secondInputs = []int{0, 0, 0, 0, 0, 0}
	gates := make([]ArithGate, inputSize*2)
	firstInputs := make([]int, inputSize*2)
	secondInputs := make([]int, inputSize*2)
	for i := 0; i < inputSize; i++ {
		gates[i] = InputA
		firstInputs[i] = i
		secondInputs[i] = 0
	}
	for i := 0; i < inputSize; i++ {
		gates[inputSize+i] = InputB
		firstInputs[inputSize+i] = i
		secondInputs[i] = 0
	}
	return gates, firstInputs, secondInputs
}

func getDotProductSubCircuit(inputSize int) ([]ArithGate, []int, []int) {
	offset := inputSize * 2
	gates := make([]ArithGate, inputSize*2)
	firstInputs := make([]int, inputSize*2)
	secondInputs := make([]int, inputSize*2)

	//x[i]*w[i]
	for i := 0; i < inputSize; i++ {
		gates[i] = Mul2Wires
		firstInputs[i] = i
		secondInputs[i] = i + inputSize
	}

	//TODO: change this to binary tree shape
	// sum of x[i]*w[i]
	gates[inputSize] = Add2Wires
	firstInputs[inputSize] = offset
	secondInputs[inputSize] = offset + 1
	for i := 1; i < inputSize; i++ {
		gates[inputSize+i] = Add2Wires
		firstInputs[inputSize+i] = offset + inputSize + i - 1
		secondInputs[inputSize+i] = offset + i + 1
	}

	return gates, firstInputs, secondInputs
}

// func getPolynomialSubCircuit(inputSize int, approxDegree int) ([]ArithGate, []int, []int, []float64) {
// 	// if approxDegree%2 == 0 {
// 	// 	approxDegree = approxDegree - 1
// 	// }
// 	offset := inputSize * 4

// 	size := (approxDegree +1)/2
// 	gates := make([]ArithGate, inputSize*2) //TODO: this size is wrong
// 	firstInputs := make([]int, inputSize*2)
// 	secondInputs := make([]int, inputSize*2)
// 	constants := make([]float64, inputSize*2)

// 	gates[0] = Mul2Wires //y^2
// 	firstInputs[0] = offset - 1
// 	secondInputs[0] = offset - 1

// 	gates[1] = Mul2Wires         //y^3
// 	firstInputs[0] = offset      //y^2
// 	secondInputs[0] = offset - 1 //y

// 	for i := 2; i < (approxDegree+1)/2; i++ {
// 		gates[i] = Mul2Wires
// 		firstInputs[i] = offset //y^2
// 		secondInputs[i] = offset + i - 1
// 	}

// 	offset2 := offset + approxDegree/2
// 	return gates, firstInputs, secondInputs
// }

func getCoefficient(degree int) float64 {

	//TODO
	return 1.0
}
