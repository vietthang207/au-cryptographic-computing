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
		if gates[i] == And2Wires {
			// ua[i] = rand.Intn(MAX_BOOL)
			// ub[i] = rand.Intn(MAX_BOOL)
			// va[i] = rand.Intn(MAX_BOOL)
			// vb[i] = rand.Intn(MAX_BOOL)
			// wa[i] = rand.Intn(MAX_BOOL)
			ua[i] = RandBigDec()
			ub[i] = RandBigDec()
			va[i] = RandBigDec()
			vb[i] = RandBigDec()
			wa[i] = RandBigDec()
			// wb[i] = wa[i] ^ ((ua[i] ^ ub[i]) & (va[i] ^ vb[i]))
			wb[i] = Add(wa[i], Mul(Add(ua[i], ub[i]), Add(va[i], vb[i])))
		}
	}
	return dealer{ua, ub, va, vb, wa, wb}
}

func (d *dealer) randA() ([]BigDec, []BigDec, []BigDec) {
	return d.ua, d.va, d.wa
}

func (d *dealer) randB() ([]BigDec, []BigDec, []BigDec) {
	return d.ub, d.vb, d.wb
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
