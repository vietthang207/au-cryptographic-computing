package main

import (
	"fmt"
	"math/rand"
)

const TABLE_SIZE = 8
const SECURITY_PARAM = 25
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

// Part 2: implement the OT protocol

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
	q, p := getSafePrimePair(SECURITY_PARAM)
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

type alice struct {
	grp   ddhGroup
	x, sk int
}

type bob struct {
	grp ddhGroup
	y   int
}

func initAlice(grp ddhGroup, x int) alice {
	sk := elGamalSampleSk(grp)
	return alice{grp: grp, x: x, sk: sk}
}

func initBob(grp ddhGroup, y int) bob {
	return bob{grp: grp, y: y}
}

func (a *alice) choose(x int) (m [TABLE_SIZE]zStarPublicKey) {
	for i := 0; i < TABLE_SIZE; i++ {
		if i == x {
			m[i] = elGamalGen(a.grp, a.sk)
		} else {
			m[i] = elGamalObliviousGen(a.grp)
		}
	}
	return
}

func (b *bob) transfer(y int, pks [TABLE_SIZE]zStarPublicKey) (m [TABLE_SIZE]zStarCipherText) {
	for i := 0; i < TABLE_SIZE; i++ {
		res := bloodTypeTruthTable(i, y)
		encodedMsg := elGamalEncode(b.grp, res)
		m[i] = elGamalEnc(b.grp, pks[i], encodedMsg)
	}
	return
}

func (a *alice) retrieve(m [TABLE_SIZE]zStarCipherText) int {
	return elGamalDecode(a.grp, elGamalDec(a.grp, a.sk, m[a.x]))
}

func simulateProtocol(grp ddhGroup, x int, y int) int {
	//TODO
	a := initAlice(grp, x)
	b := initBob(grp, y)
	m1 := a.choose(x)
	m2 := b.transfer(y, m1)
	return a.retrieve(m2)
}

func main() {
	// testBloodTypeTruthTable();

	grp := initZStarOfSafePrime()
	// Simple testing
	wrongFlag := false
	for x := 0; x < TABLE_SIZE; x++ {
		for y := 0; y < TABLE_SIZE; y++ {
			if simulateProtocol(&grp, x, y) != bloodTypeTruthTable(x, y) {
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
