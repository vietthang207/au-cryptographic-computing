package ppml

type alice struct {
	circuit     circuit
	x           []BigDec
	wires       []BigDec
	currentWire int
	ua, va, wa  []BigDec
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

func initAlice(circuit circuit, img []float64, dealer dealer) alice {
	x := FloatArrayToBigDec(img, DEFAULT_SCALAR)
	wires := make([]BigDec, len(circuit.gates))

	return alice{circuit: circuit, x: x, wires: wires, ua: dealer.ua, va: dealer.va, wa: dealer.wa}
}

func (a *alice) isSending() bool {
	switch a.circuit.gates[a.currentWire] {
	case InputA, Mul2Wires:
		return true
	}
	return false
}

func (a *alice) isReceiving() bool {
	switch a.circuit.gates[a.currentWire] {
	case InputB, Mul2Wires, Output:
		return true
	}
	return false
}

func (a *alice) handleLocalGates() {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstInput := a.circuit.firstInputs[currentWire]
	secondInput := a.circuit.secondInputs[currentWire]
	switch gate {
	case AddConst:
		a.wires[currentWire] = Add(a.wires[firstInput], IntToBigDecDefaultScalar(secondInput))
	case MulConst:
		a.wires[currentWire] = Mul(a.wires[firstInput], IntToBigDecDefaultScalar(secondInput))
	case Add2Wires:
		a.wires[currentWire] = Add(a.wires[firstInput], a.wires[secondInput])
	}
	a.currentWire++
}

func (a *alice) handleSending() []BigDec {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstInput := a.circuit.firstInputs[currentWire]
	secondInput := a.circuit.secondInputs[currentWire]
	data := make([]BigDec, 0)
	switch gate {
	case InputA:
		xb := RandForDealer()
		xa := Sub(a.x[firstInput], xb)
		a.wires[currentWire] = xa
		a.currentWire++
		return append(data, xb)
	case Mul2Wires:
		da := Add(a.wires[firstInput], a.ua[currentWire])
		ea := Add(a.wires[secondInput], a.va[currentWire])
		data = append(data, da)
		data = append(data, ea)
		return data
	default:
		panic("Incorrect case")
	}
}

func (a *alice) handleReceiving(data []BigDec) {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstInput := a.circuit.firstInputs[currentWire]
	secondInput := a.circuit.secondInputs[currentWire]
	switch gate {
	case InputB:
		a.wires[currentWire] = data[0]
		a.currentWire++
	case Mul2Wires:
		db := data[0]
		eb := data[1]
		da := Add(a.wires[firstInput], a.ua[currentWire])
		ea := Add(a.wires[secondInput], a.va[currentWire])
		d := Add(da, db)
		e := Add(ea, eb)
		a.wires[currentWire] = Sub(a.wa[currentWire], Add(Mul(e, a.wires[firstInput]), Add(Mul(d, a.wires[secondInput]), Mul(e, d))))
		a.currentWire++
	case Output:
		// fmt.Println("oa: ", a.wires[currentWire-1].integral)
		// fmt.Println("ob: ", data[0].integral)
		a.wires[currentWire] = Add(a.wires[currentWire-1], data[0])
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

func (a *alice) output() float64 {
	// ret := Mod(a.wires[a.currentWire-1], IntToBigDecDefaultScalar(2))
	ret := a.wires[a.currentWire-1]
	// fmt.Println("alice output: ", ret.ToFloat())
	return ret.ToFloat()
}
