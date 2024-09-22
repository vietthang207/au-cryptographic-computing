package main

import "fmt"
import "math/rand"

const TABLE_SIZE = 8
const MAX_BOOL = 2

// Part 1: this part generate test cases bases on the truth table
func bloodTypeTruthTable (x int, y int) int {
	if (^ (y &^ x)) & 7 == 7 {
		return 1;
	} else {
		return 0;
	}
}

func printBloodType (x int) string {
	res := "";
	switch x >> 1 {
	case 0:
		res += "O"
	case 1:
		res += "B"
	case 2:
		res += "A"
	case 3:
		res += "AB"
	}
	if x & 1 == 0 {
		res += "-" 
	} else {
		res += "+"
	}
	return res;
}

func testBloodTypeTruthTable () {
	for x := 0; x<TABLE_SIZE; x++ {
		for y := 0; y<TABLE_SIZE; y++ {
			fmt.Println(printBloodType(x), " ", printBloodType(y), " ", bloodTypeTruthTable (x, y));
		}
	}
	return;
}

// Part 2: implement the BeDOZa protocol
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
	gates []LogicGate
	firstFanins []int
	secondFanins []int
}

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

type alice struct {
	circuit circuit
	x []int
	wires []int
	currentWire int
	ua, va, wa []int
}

type bob struct {
	circuit circuit
	y []int
	wires []int
	currentWire int
	ub, vb, wb []int
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

	for i:=0; i<circuitSize; i++ {
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

func getAliceInputSize (circuit circuit) int {
	gates := circuit.gates
	circuitSize := len(gates)
	inputASize := 0
	for i:=0; i<circuitSize; i++ {
		if gates[i] == InputA {
			inputASize ++
		}
	}
	return inputASize
}

func initAlice(circuit circuit, bitmask int, dealer dealer) alice {
	aliceInputSize := getAliceInputSize(circuit)
	x := make([]int, aliceInputSize)
	for i:=0; i<aliceInputSize; i++ {
		x[aliceInputSize - i - 1] = (bitmask & (1 << i)) >> i
	}
	wires := make([]int, len(circuit.gates))

	return alice {circuit: circuit, x: x, wires: wires, ua: dealer.ua, va: dealer.va, wa: dealer.wa};
}

func (a *alice) isSending() bool {
	switch a.circuit.gates[a.currentWire] {
	case InputA, And2Wires: return true
	}
	return false
}

func (a *alice) isReceiving() bool {
	switch a.circuit.gates[a.currentWire] {
	case InputB, And2Wires, Output: return true
	}
	return false
}

func (a *alice) handleLocalGates() {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstFanin := a.circuit.firstFanins[currentWire]
	secondFanin := a.circuit.secondFanins[currentWire]
	switch gate {
	case XorConst:
		a.wires[currentWire] = a.wires[firstFanin] ^ secondFanin
	case AndConst:
		a.wires[currentWire] = a.wires[firstFanin] & secondFanin
	case Xor2Wires:
		a.wires[currentWire] = a.wires[firstFanin] ^ a.wires[secondFanin]
	}
	a.currentWire ++
	return
}

func (a *alice) handleSending() int {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstFanin := a.circuit.firstFanins[currentWire]
	secondFanin := a.circuit.secondFanins[currentWire]
	switch gate {
	case InputA:
		xb := rand.Intn(MAX_BOOL)
		xa := a.x[firstFanin] ^ xb
		a.wires[currentWire] = xa
		a.currentWire ++
		return xb
	case And2Wires:
		da := a.wires[firstFanin] ^ a.ua[currentWire]
		ea := a.wires[secondFanin] ^ a.va[currentWire]
		bitmask := da << 1 + ea
		return bitmask
	default:
		panic("Incorrect case")
	}
}

func (a *alice) handleReceiving(bitmask int) {
	currentWire := a.currentWire
	gate := a.circuit.gates[currentWire]
	firstFanin := a.circuit.firstFanins[currentWire]
	secondFanin := a.circuit.secondFanins[currentWire]
	switch gate {
	case InputB:
		a.wires[currentWire] = bitmask
		a.currentWire ++
	case And2Wires:
		db := bitmask >> 1
		eb := bitmask - db << 1
		da := a.wires[firstFanin] ^ a.ua[currentWire]
		ea := a.wires[secondFanin] ^ a.va[currentWire]
		d := da ^ db
		e := ea ^ eb
		a.wires[currentWire] = a.wa[currentWire] ^ (e & a.wires[firstFanin]) ^ (d & a.wires[secondFanin]) ^ (e & d)
		a.currentWire ++
	case Output:
		a.wires[currentWire] = a.wires[currentWire-1] ^ bitmask
		a.currentWire ++
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

func getBobInputSize (circuit circuit) int {
	gates := circuit.gates
	circuitSize := len(gates)
	inputBSize := 0
	for i:=0; i<circuitSize; i++ {
		if gates[i] == InputB {
			inputBSize ++
		}
	}
	return inputBSize
}

func initBob(circuit circuit, bitmask int, dealer dealer) bob {
	bobInputSize := getBobInputSize(circuit)

	y := make([]int, bobInputSize)
	for i:=0; i<bobInputSize; i++ {
		y[bobInputSize - i - 1] = (bitmask & (1 << i)) >> i
	}
	wires := make([]int, len(circuit.gates))
	return bob{circuit: circuit, y: y, wires: wires, ub: dealer.ub, vb: dealer.vb, wb: dealer.wb};
}

func (b *bob) isSending() bool {
	switch b.circuit.gates[b.currentWire] {
	case InputB, And2Wires, Output: return true
	}
	return false
}

func (b *bob) isReceiving() bool {
	switch b.circuit.gates[b.currentWire] {
	case InputA, And2Wires: return true
	}
	return false
}

func (b *bob) handleSending() int {
	currentWire := b.currentWire
	gate := b.circuit.gates[currentWire]
	firstFanin := b.circuit.firstFanins[currentWire]
	secondFanin := b.circuit.secondFanins[currentWire]
	switch gate {
	case InputB:
		xa := rand.Intn(MAX_BOOL)
		xb := b.y[firstFanin] ^ xa
		b.wires[currentWire] = xb
		b.currentWire ++
		return xa
	case And2Wires:
		db := b.wires[firstFanin] ^ b.ub[currentWire]
		eb := b.wires[secondFanin] ^ b.vb[currentWire]
		bitmask := db << 1 + eb
		b.currentWire ++
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
	firstFanin := b.circuit.firstFanins[currentWire]
	secondFanin := b.circuit.secondFanins[currentWire]
	switch gate {
	case InputA:
		b.wires[currentWire] = bitmask
		b.currentWire ++
	case And2Wires:
		da := bitmask >> 1
		ea := bitmask - da << 1
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
	firstFanin := b.circuit.firstFanins[currentWire]
	secondFanin := b.circuit.secondFanins[currentWire]
	switch gate {
	case XorConst:
		//Different from Alice: no ^ c
		b.wires[currentWire] = b.wires[firstFanin]
	case AndConst:
		b.wires[currentWire] = b.wires[firstFanin] & secondFanin
	case Xor2Wires:
		b.wires[currentWire] = b.wires[firstFanin] ^ b.wires[secondFanin]
	}
	b.currentWire ++
	return
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

// for debugging purpose
// func printWires(arr []int) {
// 	for i:=0; i<len(arr); i++ {
// 		fmt.Print(arr[i], " ")
// 	}
// 	fmt.Println()
// }

func simulateProtocol(circuit circuit, x int, y int, d dealer) int {
	a := initAlice(circuit, x, d)
	b := initBob(circuit, y, d)
	for !a.hasOutput() {
		receive(&b, send(&a))
		receive(&a, send(&b))
	}
	return a.output()
}

func main() {
	// testBloodTypeTruthTable();

	//Circuit encoding convention:
	//Gate                  | firstFanin          | secondFanin
	//Input                 | index of input bit  | 0
	//XOR/AND with constant | index of input wire | constant                   |
	//Binary gate           | index of input wire | index of second input wire |
	var gates = []LogicGate {InputA, InputA, InputA, InputB, InputB, InputB, XorConst, XorConst, XorConst, And2Wires, And2Wires, And2Wires, XorConst, XorConst, XorConst, And2Wires, And2Wires, Output}
	var firstFanins = []int {     0,      1,      2,      0,      1,      2,        0,        1,        2,         6,         7,         8,        9,       10,       11,        12,        15,      0}
	var secondFanins = []int{     0,      0,      0,      0,      0,      0,        1,        1,        1,         3,         4,         5,        1,        1,        1,        13,        14,      0}
	bloodTypecircuit := circuit{gates: gates, firstFanins: firstFanins, secondFanins: secondFanins}
	d := initDealer(bloodTypecircuit)
	// Simple testing
	for x:=0; x<TABLE_SIZE; x++ {
		for y:=0; y<TABLE_SIZE; y++ {
			if simulateProtocol(bloodTypecircuit, x, y, d) != bloodTypeTruthTable(x, y) {
				fmt.Println("Wrong case ", x, " ", y)
			}
		}
	}
}
