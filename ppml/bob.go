package ppml

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

func initBob(circuit circuit, weights []float64, dealer dealer) bob {
	// bobInputSize := getBobInputSize(circuit)
	y := FloatArrayToBigDec(weights, DEFAULT_SCALAR)
	wires := make([]BigDec, len(circuit.gates))
	return bob{circuit: circuit, y: y, wires: wires, ub: dealer.ub, vb: dealer.vb, wb: dealer.wb}
}

func (b *bob) isSending() bool {
	switch b.circuit.gates[b.currentWire] {
	case InputB, Mul2Wires, Output:
		return true
	}
	return false
}

func (b *bob) isReceiving() bool {
	switch b.circuit.gates[b.currentWire] {
	case InputA, Mul2Wires:
		return true
	}
	return false
}

func (b *bob) handleSending() []BigDec {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstInput := b.circuit.firstInputs[currentWire]
	secondInput := b.circuit.secondInputs[currentWire]
	data := make([]BigDec, 0)
	switch gate {
	case InputB:
		xa := RandForDealer()
		xb := Sub(b.y[firstInput], xa)
		b.wires[currentWire] = xb
		b.currentWire++
		return append(data, xa)
	case Mul2Wires:
		db := Add(b.wires[firstInput], b.ub[currentWire])
		eb := Add(b.wires[secondInput], b.vb[currentWire])
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

func (b *bob) handleReceiving(data []BigDec) {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstInput := b.circuit.firstInputs[currentWire]
	secondInput := b.circuit.secondInputs[currentWire]
	switch gate {
	case InputA:
		b.wires[currentWire] = data[0]
		b.currentWire++
	case Mul2Wires:
		da := data[0]
		ea := data[1]
		db := Add(b.wires[firstInput], b.ub[currentWire])
		eb := Add(b.wires[secondInput], b.vb[currentWire])
		d := Add(da, db)
		e := Add(ea, eb)
		//Different from Alice: no ^ (e & d)
		b.wires[currentWire] = Add(b.wb[currentWire], Add(Mul(e, b.wires[firstInput]), Mul(d, b.wires[secondInput])))
	}
}

func (b *bob) handleLocalGates() {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstInput := b.circuit.firstInputs[currentWire]
	switch gate {
	case AddConst:
		//Different from Alice: no ^ c
		b.wires[currentWire] = b.wires[firstInput]
	case MulConst:
		secondInput := b.circuit.constants[currentWire]
		b.wires[currentWire] = Mul(b.wires[firstInput], secondInput)
	case Add2Wires:
		secondInput := b.circuit.secondInputs[currentWire]
		b.wires[currentWire] = Add(b.wires[firstInput], b.wires[secondInput])
	}
	b.currentWire++
}
