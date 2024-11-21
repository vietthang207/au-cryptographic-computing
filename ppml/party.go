package ppml

type dealer struct {
	ua, ub, va, vb, wa, wb []BigDec
}

type party interface {
	isSending() bool
	isReceiving() bool
	handleSending() []BigDec
	handleReceiving([]BigDec)
	handleLocalGates()
}

func initDealer(circuit circuit) dealer {
	gates := circuit.gates
	circuitSize := len(gates)
	ua := make([]BigDec, circuitSize)
	ub := make([]BigDec, circuitSize)
	va := make([]BigDec, circuitSize)
	vb := make([]BigDec, circuitSize)
	wa := make([]BigDec, circuitSize)
	wb := make([]BigDec, circuitSize)

	for i := 0; i < circuitSize; i++ {
		if gates[i] == Mul2Wires {
			ua[i] = RandBoolBigDec()
			ub[i] = RandBoolBigDec()
			va[i] = RandBoolBigDec()
			vb[i] = RandBoolBigDec()
			wa[i] = RandBoolBigDec()
			wb[i] = Add(wa[i], Mul(Add(ua[i], ub[i]), Add(va[i], vb[i])))
		}
	}
	return dealer{ua, ub, va, vb, wa, wb}
}

func send(p party) []BigDec {
	for !p.isSending() {
		if p.isReceiving() {
			data := make([]BigDec, 0)
			return data
		}
		p.handleLocalGates()
	}
	return p.handleSending()
}

func receive(p party, data []BigDec) {
	for !p.isReceiving() {
		if p.isSending() {
			return
		}
		p.handleLocalGates()
	}
	p.handleReceiving(data)
}
