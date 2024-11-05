package main

type LogicGate int

const (
	InputA LogicGate = iota
	InputB
	Output
	XorConst
	AndConst
	Xor2Wires
	And2Wires
)

type circuit struct {
	gates        []LogicGate
	firstInputs  []int
	secondInputs []int
}
