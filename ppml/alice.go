package main

import "math/rand"

type alice struct {
	circuit     circuit
	x           []int
	wires       []int
	currentWire int
	ua, va, wa  []int
}

func getAliceInputSize(circuit circuit) int {
	gates := circuit.gates
	circuitSize := len(gates)
	inputASize := 0
	for i := 0; i < circuitSize; i++ {
		if gates[i] == InputA {
			inputASize++
		}
	}
	return inputASize
}

func initAlice(circuit circuit, bitmask int, dealer dealer) alice {
	aliceInputSize := getAliceInputSize(circuit)
	x := make([]int, aliceInputSize)
	for i := 0; i < aliceInputSize; i++ {
		x[aliceInputSize-i-1] = (bitmask & (1 << i)) >> i
	}
	wires := make([]int, len(circuit.gates))

	return alice{circuit: circuit, x: x, wires: wires, ua: dealer.ua, va: dealer.va, wa: dealer.wa}
}

func (a *alice) isSending() bool {
	switch a.circuit.gates[a.currentWire] {
	case InputA, And2Wires:
		return true
	}
	return false
}

func (a *alice) isReceiving() bool {
	switch a.circuit.gates[a.currentWire] {
	case InputB, And2Wires, Output:
		return true
	}
	return false
}

func (a *alice) handleLocalGates() {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstFanin := a.circuit.firstInputs[currentWire]
	secondFanin := a.circuit.secondInputs[currentWire]
	switch gate {
	case XorConst:
		a.wires[currentWire] = a.wires[firstFanin] ^ secondFanin
	case AndConst:
		a.wires[currentWire] = a.wires[firstFanin] & secondFanin
	case Xor2Wires:
		a.wires[currentWire] = a.wires[firstFanin] ^ a.wires[secondFanin]
	}
	a.currentWire++
	return
}

func (a *alice) handleSending() int {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstFanin := a.circuit.firstInputs[currentWire]
	secondFanin := a.circuit.secondInputs[currentWire]
	switch gate {
	case InputA:
		xb := rand.Intn(MAX_BOOL)
		xa := a.x[firstFanin] ^ xb
		a.wires[currentWire] = xa
		a.currentWire++
		return xb
	case And2Wires:
		da := a.wires[firstFanin] ^ a.ua[currentWire]
		ea := a.wires[secondFanin] ^ a.va[currentWire]
		bitmask := da<<1 + ea
		return bitmask
	default:
		panic("Incorrect case")
	}
}

func (a *alice) handleReceiving(bitmask int) {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstFanin := a.circuit.firstInputs[currentWire]
	secondFanin := a.circuit.secondInputs[currentWire]
	switch gate {
	case InputB:
		a.wires[currentWire] = bitmask
		a.currentWire++
	case And2Wires:
		db := bitmask >> 1
		eb := bitmask - db<<1
		da := a.wires[firstFanin] ^ a.ua[currentWire]
		ea := a.wires[secondFanin] ^ a.va[currentWire]
		d := da ^ db
		e := ea ^ eb
		a.wires[currentWire] = a.wa[currentWire] ^ (e & a.wires[firstFanin]) ^ (d & a.wires[secondFanin]) ^ (e & d)
		a.currentWire++
	case Output:
		a.wires[currentWire] = a.wires[currentWire-1] ^ bitmask
		a.currentWire++
	}
}

func (a *alice) hasOutput() bool {
	if a.currentWire == len(a.circuit.gates) {
		return true
	} else {
		return false
	}
}

func (a *alice) output() int {
	return a.wires[a.currentWire-1]
}
