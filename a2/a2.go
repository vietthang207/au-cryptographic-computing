package main

import "fmt"

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
	for x := 0; x<8; x++ {
		for y := 0; y<8; y++ {
			fmt.Println(printBloodType(x), " ", printBloodType(y), " ", bloodTypeTruthTable (x, y));
		}
	}
	return;
}

func main() {
	testBloodTypeTruthTable();
}
