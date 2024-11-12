package ppml

import "math/rand"

type dealer struct {
	ua, ub, va, vb, wa, wb []int
}

type party interface {
	isSending() bool
	isReceiving() bool
	handleSending() int
	handleReceiving(int)
	handleLocalGates()
}

func initDealer(circuit circuit) dealer {
	gates := circuit.gates
	circuitSize := len(gates)
	ua := make([]int, circuitSize)
	ub := make([]int, circuitSize)
	va := make([]int, circuitSize)
	vb := make([]int, circuitSize)
	wa := make([]int, circuitSize)
	wb := make([]int, circuitSize)

	for i := 0; i < circuitSize; i++ {
		if gates[i] == And2Wires {
			ua[i] = rand.Intn(MAX_BOOL)
			ub[i] = rand.Intn(MAX_BOOL)
			va[i] = rand.Intn(MAX_BOOL)
			vb[i] = rand.Intn(MAX_BOOL)
			wa[i] = rand.Intn(MAX_BOOL)
			wb[i] = wa[i] ^ ((ua[i] ^ ub[i]) & (va[i] ^ vb[i]))
		}
	}
	return dealer{ua, ub, va, vb, wa, wb}
}

func (d *dealer) randA() ([]int, []int, []int) {
	return d.ua, d.va, d.wa
}

func (d *dealer) randB() ([]int, []int, []int) {
	return d.ub, d.vb, d.wb
}

func send(p party) int {
	for !p.isSending() {
		if p.isReceiving() {
			return 0
		}
		p.handleLocalGates()
	}
	return p.handleSending()
}

func receive(p party, bitmask int) {
	for !p.isReceiving() {
		if p.isSending() {
			return
		}
		p.handleLocalGates()
	}
	p.handleReceiving(bitmask)
}
