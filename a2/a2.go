package main

import "fmt"
import "math/rand"

const TABLE_SIZE = 8

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

func initTruthTable () [TABLE_SIZE][TABLE_SIZE]int {
	T := [TABLE_SIZE][TABLE_SIZE]int{};
	for x := 0; x<TABLE_SIZE; x++ {
		for y := 0; y<TABLE_SIZE; y++ {
			T[x][y] = bloodTypeTruthTable(x, y);
		}
	}
	return T;
}

type dealer struct {
	r, s int;
	T, Ma, Mb [TABLE_SIZE][TABLE_SIZE]int;
}

type alice struct {
	x, r, u, v, zb int;
	Ma [TABLE_SIZE][TABLE_SIZE]int;
}

type bob struct {
	y, s, u, v, zb int;
	Mb [TABLE_SIZE][TABLE_SIZE]int;
}

func initDealer() dealer {
	T := initTruthTable()
	r := rand.Intn(TABLE_SIZE)
	s := rand.Intn(TABLE_SIZE)

	var Mb [TABLE_SIZE][TABLE_SIZE]int
	for i:=0; i<TABLE_SIZE; i++ {
		for j:=0; j<TABLE_SIZE; j++ {
			Mb[i][j] = rand.Intn(2);
		}
	}

	var Ma [TABLE_SIZE][TABLE_SIZE]int
	for i:=0; i<TABLE_SIZE; i++ {
		for j:=0; j<TABLE_SIZE; j++ {
			Ma[i][j] = Mb[i][j] ^ T[(i-r+TABLE_SIZE)%TABLE_SIZE][(j-s+TABLE_SIZE)%TABLE_SIZE]
		}
	}
	return dealer{T: T, Ma: Ma, Mb: Mb, r: r, s: s}
}

func (d *dealer) randA() (int, [TABLE_SIZE][TABLE_SIZE]int) {
	return d.r, d.Ma
}

func (d *dealer) randB() (int, [TABLE_SIZE][TABLE_SIZE]int) {
	return d.s, d.Mb
}

func initAlice(x int, r int, Ma [TABLE_SIZE][TABLE_SIZE]int) alice {
	return alice {x: x, r: r, Ma: Ma};
}

func (a *alice) send() int {
	a.u = (a.x + a.r) % TABLE_SIZE
	return a.u
}

func (a *alice) receive(v int, zb int) {
	a.v = v
	a.zb = zb
}

func (a *alice) output() int {
	return a.Ma[a.u][a.v] ^ a.zb
}

func initBob(y int, s int, Mb [TABLE_SIZE][TABLE_SIZE]int) bob {
	return bob {y: y, s: s, Mb: Mb};
}

func (b *bob) receive(u int) (int, int) {
	b.u = u
	b.v = (b.y + b.s) % TABLE_SIZE
	b.zb = b.Mb[b.u][b.v]
	return b.v, b.zb
}

func (b *bob) send() (int, int) {
	return b.v, b.zb
}

func simulateProtocol(x int, y int, d dealer) int {
	r, Ma := d.randA()
	a := initAlice(x, r, Ma)
	s, Mb := d.randB()
	b := initBob(y, s, Mb)
	b.receive(a.send())

	v, zb := b.send()
	a.receive(v, zb)
	return a.output()
}

func main() {
	// testBloodTypeTruthTable();
	d := initDealer()
	// Simple testing
	for x:=0; x<TABLE_SIZE; x++ {
		for y:=0; y<TABLE_SIZE; y++ {
			if simulateProtocol(x, y, d) != bloodTypeTruthTable(x, y) {
				fmt.Println("Wrong case ", x, " ", y)
			}
		}
	}
}
