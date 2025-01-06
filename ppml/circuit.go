package ppml

import "fmt"

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
	g1, f1, s1 := getInputSubCircuit(inputSize)
	g2, f2, s2 := getDotProductSubCircuit(inputSize)
	g3, f3, s3, c3 := getPolynomialSubCircuit(inputSize, approxDegree)

	g := append(append(g1, g2...), g3...)
	f := append(append(f1, f2...), f3...)
	s := append(append(s1, s2...), s3...)

	constFloats := make([]float64, inputSize*4)
	constFloats = append(constFloats, c3...)
	constants := FloatArrayToBigDec(constFloats, DEFAULT_SCALAR)

	g = append(g, Output)
	f = append(f, 0)
	s = append(s, 0)
	return circuit{gates: g, firstInputs: f, secondInputs: s, constants: constants}
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

func getPolynomialSubCircuit(inputSize int, approxDegree int) ([]ArithGate, []int, []int, []float64) {
	offset := inputSize * 4

	size := 0
	if approxDegree%2 == 0 {
		size = approxDegree / 2
	} else {
		size = (approxDegree + 1) / 2
	}

	gates := make([]ArithGate, size*3)
	firstInputs := make([]int, size*3)
	secondInputs := make([]int, size*3)
	constants := make([]float64, size*3)

	gates[0] = Mul2Wires //y^2
	firstInputs[0] = offset - 1
	secondInputs[0] = offset - 1

	gates[1] = Mul2Wires         //y^3
	firstInputs[1] = offset      //y^2
	secondInputs[1] = offset - 1 //y

	for i := 2; i < size; i++ {
		gates[i] = Mul2Wires
		firstInputs[i] = offset //y^2
		secondInputs[i] = offset + i - 1
	}

	gates[size] = MulConst         //a1 * y
	firstInputs[size] = offset - 1 // y
	constants[size] = getCoefficient(1)

	for i := 1; i < size; i++ {
		// a_(2i+1) * y^(2i+1)
		gates[size+i] = MulConst
		firstInputs[size+i] = offset + i
		constants[size+i] = getCoefficient(2*i + 1)
	}

	gates[2*size] = AddConst
	firstInputs[2*size] = offset + size
	constants[2*size] = 0.5

	for i := 1; i < size; i++ {
		gates[2*size+i] = Add2Wires
		firstInputs[2*size+i] = offset + size + i
		secondInputs[2*size+i] = offset + 2*size + i - 1
	}

	return gates, firstInputs, secondInputs, constants
}

// func fact(n int) float64 {
// 	if n == 0 {
// 		return 1
// 	}
// 	return float64(n) * fact(n-1)
// }

// func bernoulli(n int) float64 {
// 	B := make([]float64, n+1)
// 	A := make([]float64, n+1)

// 	for m := 0; m <= n; m++ {
// 		A[m] = 1.0 / float64(m+1)
// 		for j := m; j > 0; j-- {
// 			A[j-1] = float64(j) * (A[j-1] - A[j])
// 		}
// 		B[m] = A[0]
// 	}

// 	return B[n]
// }

// func getCoefficient(degree int) float64 {
// 	if degree == 0 {
// 		return 0.5
// 	}
// 	if degree == 1 {
// 		return 0.25
// 	}
// 	if degree%2 == 0 {
// 		return 0.0
// 	}

// 	tmp := (1<<(degree-1) - 1)
// 	res := 0.25 * float64(tmp) * bernoulli(degree-1) / fact(degree)
// 	if (degree-1)%2 == 1 {
// 		res = -res
// 	}
// 	return res
// }

func getCoefficient(degree int) float64 {
	coef := [20]float64{0.5, 0.25, 0, -1.0 / 48, 0, 1.0 / 480, 0, -17.0 / 80640, 0, 31.0 / 1451520, 0, -0.0000021645, 0, 0.00000020498, 0, -0.0000000198, 0, 0.00000001719, 0, -0.00000001191}
	return coef[degree]
}

func TestCoefficient() {
	for i := 0; i < 10; i++ {
		fmt.Println(getCoefficient(i))
	}
}
