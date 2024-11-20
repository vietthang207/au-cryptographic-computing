package ppml

import "math/rand"

type bob struct {
	circuit     circuit
	y           []BigDec
	wires       []BigDec
	currentWire int
	ub, vb, wb  []BigDec
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

	y := make([]BigDec, bobInputSize)
	for i := 0; i < bobInputSize; i++ {
		y[bobInputSize-i-1] = IntToBigDecDefaultScalar((bitmask & (1 << i)) >> i)
	}
	wires := make([]BigDec, len(circuit.gates))
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

func (b *bob) handleSending() []BigDec {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstFanin := b.circuit.firstInputs[currentWire]
	secondFanin := b.circuit.secondInputs[currentWire]
	data := make([]int, 0)
	switch gate {
	case InputB:
		xa := rand.Intn(MAX_BOOL)
		xb := b.y[firstFanin] + xa
		b.wires[currentWire] = xb
		b.currentWire++
		return append(data, xa)
	case And2Wires:
		db := b.wires[firstFanin] + b.ub[currentWire]
		eb := b.wires[secondFanin] + b.vb[currentWire]
		//bitmask := db<<1 + eb
		data = append(data, db)
		data = append(data, eb)
		b.currentWire++
		return data
	case Output:
		openValue := b.wires[currentWire-1]
		return append(data, openValue)
	default:
		panic("Incorrect case")
	}
}

func (b *bob) handleReceiving(data []int) {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstFanin := b.circuit.firstInputs[currentWire]
	secondFanin := b.circuit.secondInputs[currentWire]
	switch gate {
	case InputA:
		b.wires[currentWire] = data[0]
		b.currentWire++
	case And2Wires:
		da := data[0]
		ea := data[1]
		db := b.wires[firstFanin] + b.ub[currentWire]
		eb := b.wires[secondFanin] + b.vb[currentWire]
		d := da + db
		e := ea + eb
		//Different from Alice: no ^ (e & d)
		b.wires[currentWire] = b.wb[currentWire] + (e * b.wires[firstFanin]) + (d * b.wires[secondFanin])
	}
}

func (b *bob) handleLocalGates() {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstFanin := b.circuit.firstInputs[currentWire]
	secondFanin := b.circuit.secondInputs[currentWire]
	switch gate {
	case AddConst:
		//Different from Alice: no ^ c
		b.wires[currentWire] = b.wires[firstFanin]
	case AndConst:
		b.wires[currentWire] = Mul(b.wires[firstFanin], IntToBigDecDefaultScalar(secondFanin))
	case Xor2Wires:
		b.wires[currentWire] = Add(b.wires[firstFanin], b.wires[secondFanin])
	}
	b.currentWire++
	return
}
