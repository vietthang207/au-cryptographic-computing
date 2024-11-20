package ppml

type LogicGate int

const (
	InputA LogicGate = iota
	InputB
	Output
	AddConst
	AndConst
	Xor2Wires
	And2Wires
)

type circuit struct {
	gates        []LogicGate
	firstInputs  []int
	secondInputs []int
}
