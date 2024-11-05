package main

import "math/rand"

type bob struct {
	circuit     circuit
	y           []int
	wires       []int
	currentWire int
	ub, vb, wb  []int
}

func getBobInputSize(circuit circuit) int {
	gates := circuit.gates
	circuitSize := len(gates)
	inputBSize := 0
	for i := 0; i < circuitSize; i++ {
		if gates[i] == InputB {
			inputBSize++
		}
	}
	return inputBSize
}

func initBob(circuit circuit, bitmask int, dealer dealer) bob {
	bobInputSize := getBobInputSize(circuit)

	y := make([]int, bobInputSize)
	for i := 0; i < bobInputSize; i++ {
		y[bobInputSize-i-1] = (bitmask & (1 << i)) >> i
	}
	wires := make([]int, len(circuit.gates))
	return bob{circuit: circuit, y: y, wires: wires, ub: dealer.ub, vb: dealer.vb, wb: dealer.wb}
}

func (b *bob) isSending() bool {
	switch b.circuit.gates[b.currentWire] {
	case InputB, And2Wires, Output:
		return true
	}
	return false
}

func (b *bob) isReceiving() bool {
	switch b.circuit.gates[b.currentWire] {
	case InputA, And2Wires:
		return true
	}
	return false
}

func (b *bob) handleSending() int {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstFanin := b.circuit.firstInputs[currentWire]
	secondFanin := b.circuit.secondInputs[currentWire]
	switch gate {
	case InputB:
		xa := rand.Intn(MAX_BOOL)
		xb := b.y[firstFanin] ^ xa
		b.wires[currentWire] = xb
		b.currentWire++
		return xa
	case And2Wires:
		db := b.wires[firstFanin] ^ b.ub[currentWire]
		eb := b.wires[secondFanin] ^ b.vb[currentWire]
		bitmask := db<<1 + eb
		b.currentWire++
		return bitmask
	case Output:
		openValue := b.wires[currentWire-1]
		return openValue
	default:
		panic("Incorrect case")
	}
}

func (b *bob) handleReceiving(bitmask int) {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstFanin := b.circuit.firstInputs[currentWire]
	secondFanin := b.circuit.secondInputs[currentWire]
	switch gate {
	case InputA:
		b.wires[currentWire] = bitmask
		b.currentWire++
	case And2Wires:
		da := bitmask >> 1
		ea := bitmask - da<<1
		db := b.wires[firstFanin] ^ b.ub[currentWire]
		eb := b.wires[secondFanin] ^ b.vb[currentWire]
		d := da ^ db
		e := ea ^ eb
		//Different from Alice: no ^ (e & d)
		b.wires[currentWire] = b.wb[currentWire] ^ (e & b.wires[firstFanin]) ^ (d & b.wires[secondFanin])
	}
}

func (b *bob) handleLocalGates() {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstFanin := b.circuit.firstInputs[currentWire]
	secondFanin := b.circuit.secondInputs[currentWire]
	switch gate {
	case XorConst:
		//Different from Alice: no ^ c
		b.wires[currentWire] = b.wires[firstFanin]
	case AndConst:
		b.wires[currentWire] = b.wires[firstFanin] & secondFanin
	case Xor2Wires:
		b.wires[currentWire] = b.wires[firstFanin] ^ b.wires[secondFanin]
	}
	b.currentWire++
	return
}
