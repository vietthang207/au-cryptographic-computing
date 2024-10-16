package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
)

const DDH_PARAM_K = 26
const PARAM_K = 15
const PARAM_N = 3
const TABLE_SIZE = 1 << PARAM_N
const MAX_RAND = 1 << 15

// Part 1: this part generate test cases bases on the truth table
func bloodTypeTruthTable(x int, y int) int {
	if (^(y &^ x))&7 == 7 {
		return 1
	} else {
		return 0
	}
}

func printBloodType(x int) string {
	res := ""
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
	if x&1 == 0 {
		res += "-"
	} else {
		res += "+"
	}
	return res
}

func testBloodTypeTruthTable() {
	for x := 0; x < TABLE_SIZE; x++ {
		for y := 0; y < TABLE_SIZE; y++ {
			fmt.Println(printBloodType(x), " ", printBloodType(y), " ", bloodTypeTruthTable(x, y))
		}
	}
	return
}

// Part 2: implement the Yao's protocol

type LogicGate int

const (
	InputA LogicGate = iota
	InputB
	Not
	Xor
	And
)

type circuit struct {
	gates        []LogicGate
	firstFanins  []int
	secondFanins []int
}

// Fast exponentiation using divide and conquer
func fastExpMod(a int, b int, mod int) int {
	if b == 0 {
		return 1
	}
	tmp := fastExpMod(a, b/2, mod)
	if b%2 == 0 {
		return (tmp * tmp) % mod
	} else {
		return ((tmp * tmp) % mod) % mod
	}
}

func extendedGcd(a int, b int) (int, int, int) {
	prev_x, x := 1, 0
	prev_y, y := 0, 1
	for b != 0 {
		q := a / b
		prev_x, x = x, prev_x-q*x
		prev_y, y = y, prev_y-q*y
		a, b = b, a%b
	}
	return a, prev_x, prev_y
}

func isPrime(x int) bool {
	//TODO: improve this using Miller-Rabin
	for i := 2; i < x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func getSafePrimePair(size int) (q int, p int) {
	min := 1 << (size - 1)
	// -2 and +1 to ensure that we do not output 1 or -1
	r := rand.Intn(min) - 2
	q = min + r + 1
	p = 2*q + 1
	for !(isPrime(p) && isPrime(q)) {
		r := rand.Intn(min) - 2
		q = min + r + 1
		p = 2*q + 1
	}
	fmt.Println("Safe prime: q: ", q, " p: ", p)
	return
}

type ddhGroup interface {
	getOrder() int
	getGenerator() int
	multiply(int, int) int
	inverse(int) int
	pow(int, int) int
}

type zStarOfSafePrime struct {
	q, g int
}

type zStarPublicKey struct {
	g, h int
}

type zStarCipherText struct {
	c1, c2 int
}

func initZStarOfSafePrime() zStarOfSafePrime {
	q, p := getSafePrimePair(DDH_PARAM_K)
	s := rand.Intn(p)
	g := (s * s) % p
	return zStarOfSafePrime{q: q, g: g}
}

func (grp *zStarOfSafePrime) getOrder() int {
	return grp.q
}

func (grp *zStarOfSafePrime) getGenerator() int {
	return grp.g
}

func (grp *zStarOfSafePrime) multiply(a int, b int) int {
	q := grp.getOrder()
	p := 2*q + 1
	return (a * b) % p
}

func (grp *zStarOfSafePrime) inverse(a int) int {
	q := grp.getOrder()
	p := 2*q + 1
	_, tmp, _ := extendedGcd(a, p)
	return tmp % p
}
func (grp *zStarOfSafePrime) pow(a int, b int) int {
	q := grp.getOrder()
	p := 2*q + 1
	return fastExpMod(a, b, p)
}

func elGamalGen(grp ddhGroup, sk int) zStarPublicKey {
	g := grp.getGenerator()
	h := grp.pow(g, sk)
	return zStarPublicKey{g: g, h: h}
}

func elGamalEnc(grp ddhGroup, pk zStarPublicKey, m int) zStarCipherText {
	q := grp.getOrder()
	r := rand.Intn(q)
	c1 := grp.pow(pk.g, r)
	c2 := grp.multiply(m, grp.pow(pk.h, r))
	return zStarCipherText{c1: c1, c2: c2}
}

func elGamalDec(grp ddhGroup, sk int, c zStarCipherText) (m int) {
	m = grp.multiply(c.c2, grp.pow(c.c1, grp.inverse(sk)))
	return
}

func elGamalSampleSk(grp ddhGroup) (alpha int) {
	alpha = rand.Intn(grp.getOrder())
	return
}

func elGamalObliviousGen(grp ddhGroup) zStarPublicKey {
	q := grp.getOrder()
	p := 2*q + 1
	g := grp.getGenerator()
	s := rand.Intn(p)
	h := (s * s) % p
	return zStarPublicKey{g: g, h: h}
}

func elGamalEncode(grp ddhGroup, m int) int {
	q := grp.getOrder()
	p := 2*q + 1
	if grp.pow(m+1, q) == 1 {
		return (m + 1) % p
	} else {
		return (p - m - 1) % p
	}
}

func elGamalDecode(grp ddhGroup, c int) int {
	q := grp.getOrder()
	p := 2*q + 1
	if c <= q {
		return (p + c - 1) % p
	} else {
		return (p - c - 1) % p
	}
}

func getIthBit(x int, i int) int {
	return (x & (1 << i)) >> i
}

type alice struct {
	circuit circuit
	grp     ddhGroup
	x       int
	sk      [PARAM_N]int
}

type bob struct {
	circuit    circuit
	grp        ddhGroup
	y          int
	k          [][2]int
	garbTables [][4]int
}

func initAlice(circuit circuit, grp ddhGroup, x int) alice {
	var sk [PARAM_N]int
	for i := 0; i < PARAM_N; i++ {
		sk[i] = elGamalSampleSk(grp)
	}
	return alice{circuit: circuit, grp: grp, x: x, sk: sk}
}

func prf(kLeft int, kRight int, i int) int {
	//TODO
	k := kLeft<<PARAM_K + kRight
	bytes := make([]byte, 32)
	binary.LittleEndian.PutUint32(bytes, uint32(k^i))
	digest := sha256.Sum256(bytes)
	return int(binary.LittleEndian.Uint32(digest[:2*PARAM_K]))
}

func makeGarbTable(gate LogicGate, i int, k [2]int, kLeft [2]int, kRight [2]int) [4]int {
	var c, cPerm [4]int
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			c[a<<1+b] = prf(kLeft[a], kRight[b], i) ^ (k[evalGate(gate, a, b)] << PARAM_K)
		}
	}

	perm := rand.Perm(len(c))
	for i, v := range perm {
		cPerm[v] = c[i]
	}
	return cPerm
}

func initBob(circuit circuit, grp ddhGroup, y int) bob {
	circuitSize := len(circuit.gates)
	k := make([][2]int, circuitSize)
	for i := 0; i < circuitSize; i++ {
		k[i][0] = rand.Intn(1 << PARAM_K)
		k[i][1] = rand.Intn(1 << PARAM_K)
	}
	garbTables := make([][4]int, circuitSize)
	for i := 0; i < circuitSize; i++ {
		gate := circuit.gates[i]
		left := circuit.firstFanins[i]
		right := circuit.secondFanins[i]
		if gate == Not || gate == Xor || gate == And {
			garbTables[i] = makeGarbTable(gate, i, k[i], k[left], k[right])
		}
	}
	return bob{circuit: circuit, grp: grp, y: y, k: k, garbTables: garbTables}
}

func (a *alice) choose(x int) (pks [PARAM_N][2]zStarPublicKey) {
	for i := 0; i < PARAM_N; i++ {
		xi := getIthBit(x, i)
		pks[i][xi] = elGamalGen(a.grp, a.sk[i])
		pks[i][1^xi] = elGamalObliviousGen(a.grp)
	}
	return
}

func garbEncode(ex [PARAM_N][2]int, x int) [PARAM_N]int {
	var garbX [PARAM_N]int
	for i := 0; i < PARAM_N; i++ {
		xi := getIthBit(x, i)
		garbX[i] = ex[i][xi]
	}
	return garbX
}

func evalGate(gate LogicGate, firstFanin int, secondFanin int) int {
	switch gate {
	case Not:
		return firstFanin ^ 1
	case Xor:
		return firstFanin ^ secondFanin
	case And:
		return firstFanin & secondFanin
	}
	panic("unreachable case")
}

func (b *bob) transfer(y int, pks [PARAM_N][2]zStarPublicKey) (otMessages [PARAM_N][2]zStarCipherText, garbTables [][4]int, garbY [PARAM_N]int, d [2]int) {
	circuitSize := len(b.circuit.gates)
	ex := b.k[:PARAM_N]
	for i := 0; i < PARAM_N; i++ {
		encodedMsg := elGamalEncode(b.grp, ex[i][0])
		otMessages[i][0] = elGamalEnc(b.grp, pks[i][0], encodedMsg)
		encodedMsg = elGamalEncode(b.grp, ex[i][1])
		otMessages[i][1] = elGamalEnc(b.grp, pks[i][1], encodedMsg)
	}

	garbTables = b.garbTables

	var ey [PARAM_N][2]int
	for i := 0; i < PARAM_N; i++ {
		ey[i] = b.k[PARAM_N+i]
	}
	garbY = garbEncode(ey, y)

	d[0] = b.k[circuitSize-1][0]
	d[1] = b.k[circuitSize-1][1]

	return
}

func garbEval(garbTable [4]int, kLeft int, kRight int, i int) int {
	//TODO
	var ret int
	matchCount := 0
	for j := 0; j < 4; j++ {
		tmp := prf(kLeft, kRight, i) ^ garbTable[j]
		mask := 1<<PARAM_K - 1
		if tmp&mask == 0 {
			ret = tmp >> PARAM_K
			matchCount++
		}
	}
	if matchCount == 1 {
		return ret
	} else if matchCount > 1 {
		panic("more than 1 match")
	}
	panic("unable to garbEval")
}

func (a *alice) retrieve(otMessages [PARAM_N][2]zStarCipherText, garbTables [][4]int, garbY [PARAM_N]int, d [2]int) int {
	var garbX [PARAM_N]int
	for i := 0; i < PARAM_N; i++ {
		xi := getIthBit(a.x, i)
		garbX[i] = elGamalDecode(a.grp, elGamalDec(a.grp, a.sk[i], otMessages[i][xi]))
	}
	garbWires := make([]int, len(a.circuit.gates))
	for i := 0; i < PARAM_N; i++ {
		garbWires[i] = garbX[i]
		garbWires[i+PARAM_N] = garbY[i]
	}
	for i := 2 * PARAM_N; i < len(garbWires); i++ {
		left := a.circuit.firstFanins[i]
		right := a.circuit.secondFanins[i]
		garbWires[i] = garbEval(garbTables[i], garbWires[left], garbWires[right], i)
	}

	z := garbWires[len(garbWires)-1]
	if z == d[0] {
		return 0
	} else if z == d[1] {
		return 1
	}
	panic("unable to decode")
}

func simulateProtocol(circuit circuit, grp ddhGroup, x int, y int) int {
	a := initAlice(circuit, grp, x)
	b := initBob(circuit, grp, y)
	m1 := a.choose(x)
	m2, m3, m4, m5 := b.transfer(y, m1)
	return a.retrieve(m2, m3, m4, m5)
}

func main() {
	// testBloodTypeTruthTable();

	grp := initZStarOfSafePrime()

	var gates = []LogicGate{InputA, InputA, InputA, InputB, InputB, InputB, Not, Not, Not, And, And, And, Not, Not, Not, And, And}
	var firstFanins = []int{0, 1, 2, 0, 1, 2, 0, 1, 2, 6, 7, 8, 9, 10, 11, 12, 15}
	var secondFanins = []int{0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 9, 10, 11, 13, 14}
	bloodTypecircuit := circuit{gates: gates, firstFanins: firstFanins, secondFanins: secondFanins}
	// Simple testing
	wrongFlag := false
	for x := 0; x < TABLE_SIZE; x++ {
		for y := 0; y < TABLE_SIZE; y++ {
			if simulateProtocol(bloodTypecircuit, &grp, x, y) != bloodTypeTruthTable(x, y) {
				fmt.Println("Wrong case ", x, " ", y)
				wrongFlag = true
			}
		}
	}
	if wrongFlag {
		fmt.Println("There's some wrong cases")
	} else {
		fmt.Println("All cases are correct")
	}
}
