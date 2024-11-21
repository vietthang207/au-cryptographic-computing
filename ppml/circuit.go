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
}
