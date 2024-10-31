package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"os"
)

const INPUT_LENGTH = 3
const TABLE_SIZE = 1 << INPUT_LENGTH

// bit length of secret odd number p
const η = 15

// Number of noise r and quotent q_i
const N = 4

// bit length of each quotent integer q_i
const γ = 5

// bit length of noise r_i
const ρ = 2

// Part 1: this part generate test cases bases on tdhe truth table
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

// Part 2: implement d-HE
func getIthBit(x uint64, i int) uint64 {
	return (x & (1 << i)) >> i
}

func randWithBitlength(bitLength int) uint64 {
	tmp := rand.Int63n(1 << (bitLength - 1))
	return uint64(tmp + 1<<(bitLength-1))
}

func randSubsetMask() uint64 {
	subsetMask := rand.Int63n(1 << N)
	bitcount := bits.OnesCount64(uint64(subsetMask))
	for bitcount < N/4 || bitcount > 3*N/4 {
		subsetMask = rand.Int63n(1 << N)
		bitcount = bits.OnesCount64(uint64(subsetMask))
	}
	return uint64(subsetMask)
}

func dheKeygen() (sk uint64, pk [N]uint64) {
	sk = randWithBitlength(η)
	for sk%2 == 0 {
		sk = randWithBitlength(η)
	}
	fmt.Println("sk", sk)
	for i := 0; i < N; i++ {
		qi := randWithBitlength(γ)
		ri := randWithBitlength(ρ)
		pk[i] = sk*qi + 2*ri
		fmt.Println("qi: ", qi, " ri: ", ri, " yi: ", pk[i])
	}
	return
}

func dheEnc(m uint64, pk [N]uint64) uint64 {
	subsetMask := randSubsetMask()
	sum := uint64(0)
	for i := 0; i < N; i++ {
		sum += pk[i] * getIthBit(subsetMask, i)
	}
	// fmt.Println("plain: ", m)
	// fmt.Println("enc: ", m+sum)
	return m + sum
}

func dheDec(c uint64, sk uint64) uint64 {
	return (c % sk) % 2
}

type alice struct {
	sk uint64
	pk [N]uint64
}

type bob struct {
}

func initAlice(sk uint64, pk [N]uint64) alice {
	return alice{sk: sk, pk: pk}
}

func initBob() bob {
	return bob{}
}

func (a *alice) choose(x uint64) (encryptedX [INPUT_LENGTH]uint64) {
	for i := 0; i < INPUT_LENGTH; i++ {
		xi := getIthBit(x, i)
		// neXi := (1 + xi) % 2
		encryptedX[i] = dheEnc(xi, a.pk)
	}
	return
}

func eval(x [INPUT_LENGTH]uint64, y [INPUT_LENGTH]uint64) (res uint64) {
	res = 1
	for i := 0; i < INPUT_LENGTH; i++ {
		// fmt.Println(x[i], " ", y[i], " ")
		res *= (1 + (1+x[i])*y[i])
	}
	// fmt.Println("\nres: ", res)
	return
}

func (b *bob) transfer(y uint64, encryptedX [INPUT_LENGTH]uint64) (encryptedOutput uint64) {
	var yArr [INPUT_LENGTH]uint64
	for i := 0; i < INPUT_LENGTH; i++ {
		yArr[i] = getIthBit(y, i)
	}
	return eval(encryptedX, yArr)
}

func (a *alice) retrieve(encryptedOutput uint64) uint64 {
	return dheDec(encryptedOutput, a.sk)
}

func simulateProtocol(x int, y int, sk uint64, pk [N]uint64) int {
	a := initAlice(sk, pk)
	b := initBob()
	m1 := a.choose(uint64(x))
	m2 := b.transfer(uint64(y), m1)
	return int(a.retrieve(m2))
}

func main() {
	// testBloodTypeTruthTable()

	var mr, M, minP uint64
	mr = N * (1 << ρ)
	fmt.Println("Max input noise: ", 2*mr)
	M = 2 * (mr*mr*mr + 3*mr*mr + 3*mr)
	fmt.Println("Max output noise: ", M)
	minP = 1 << (η - 1)
	fmt.Println("Min p: ", minP)
	sk, pk := dheKeygen()
	// Simple testing
	wrongFlag := false
	for x := 0; x < TABLE_SIZE; x++ {
		for y := 0; y < TABLE_SIZE; y++ {
			expected := bloodTypeTruthTable(x, y)
			actual := simulateProtocol(x, y, sk, pk)
			if expected != actual {
				fmt.Println("Wrong case ", x, " ", y)
				fmt.Println(expected)
				fmt.Println(actual)
				wrongFlag = true
				break
			}
		}
		if wrongFlag {
			break
		}
	}
	if wrongFlag {
		fmt.Println("There's some wrong cases")
		os.Exit(1)
	} else {
		fmt.Println("All cases are correct")
	}
}
