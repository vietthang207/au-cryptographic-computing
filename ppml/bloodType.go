package ppml

import "fmt"

const TABLE_SIZE = 8

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
}

func simulateProtocol(circuit circuit, x int, y int, d dealer) int {
	a := initAlice(circuit, x, d)
	b := initBob(circuit, y, d)
	for !a.hasOutput() {
		receive(&b, send(&a))
		receive(&a, send(&b))
	}
	return a.output()
}

func Main() {
	// testBloodTypeTruthTable()

	//Circuit encoding convention:
	//Gate                  | firstInput          | secondInput
	//Input                 | index of input bit  | 0
	//XOR/AND with constant | index of input wire | constant                   |
	//Binary gate           | index of input wire | index of second input wire |
	var gates = []ArithGate{InputA, InputA, InputA, InputB, InputB, InputB, AddConst, AddConst, AddConst, Mul2Wires, Mul2Wires, Mul2Wires, AddConst, AddConst, AddConst, Mul2Wires, Mul2Wires, Output}
	var firstInputs = []int{0, 1, 2, 0, 1, 2, 0, 1, 2, 6, 7, 8, 9, 10, 11, 12, 15, 0}
	var secondInputs = []int{0, 0, 0, 0, 0, 0, 1, 1, 1, 3, 4, 5, 1, 1, 1, 13, 14, 0}
	bloodTypecircuit := circuit{gates: gates, firstInputs: firstInputs, secondInputs: secondInputs}
	d := initDealer(bloodTypecircuit)
	// Simple testing
	errorFlag := false
	for x := 0; x < TABLE_SIZE; x++ {
		for y := 0; y < TABLE_SIZE; y++ {
			expected := bloodTypeTruthTable(x, y)
			actual := simulateProtocol(bloodTypecircuit, x, y, d)
			if actual != expected {
				errorFlag = true
				fmt.Println("Wrong case ", x, " ", y)
				fmt.Printf("Actual: %d, Expected: %d\n", actual, expected)
			}
		}
	}
	if errorFlag {
		fmt.Println("There are some errors.")
	} else {
		fmt.Println("No error.")
	}
}
