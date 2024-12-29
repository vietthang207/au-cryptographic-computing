package ppml

import (
	"math/big"
	"math/rand"
)

const MAX_RAND = 128
const DEFAULT_SCALAR = 10000

type BigDec struct {
	integral *big.Int
	// We use the convention that "scalar" is the inverse of the scalar, so scalar=100 represents a 1/100th scale
	// - meaning that we divide by the scalar instead of multiplying when the representation is needed
	scalar *big.Int
}

func CheckScalarDifference(x BigDec, y BigDec) {
	if x.scalar.Cmp(y.scalar) != 0 {
		panic("Different scalars not implemented")
	}
}

func Add(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	resIntegral := new(big.Int)
	resIntegral.Add(x.integral, y.integral)
	return BigDec{integral: resIntegral, scalar: x.scalar}
}

// func Sub(x BigDec, y BigDec) BigDec {
// 	CheckScalarDifference(x, y)
// 	return BigDec{integral: x.integral - y.integral, scalar: x.scalar}
// }

func Mul(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	resIntegral := new(big.Int)
	resIntegral.Mul(x.integral, y.integral)
	resIntegral.Div(resIntegral, x.scalar)
	return BigDec{integral: resIntegral, scalar: x.scalar}
}

// func Div(x BigDec, y BigDec) BigDec {
// 	CheckScalarDifference(x, y)
// 	return BigDec{integral: x.integral / y.integral * x.scalar, scalar: x.scalar}
// }

// func Mod(x BigDec, y BigDec) BigDec {
// 	CheckScalarDifference(x, y)
// 	return BigDec{integral: (x.integral / x.scalar) % (y.integral / y.scalar) * x.scalar, scalar: x.scalar}
// }

func RandBoolBigDec() BigDec {
	// return BigDec{integral: int64(rand.Intn(MAX_RAND)) * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
	return BigDec{integral: big.NewInt(rand.Int63n(MAX_RAND)), scalar: big.NewInt(DEFAULT_SCALAR)}
}

// func One() BigDec {
// 	return BigDec{integral: 1 * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
// }

// func Zero() BigDec {
// 	return BigDec{integral: 0 * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
// }

func IntToBigDecDefaultScalar(x int) BigDec {
	scalar := big.NewInt(DEFAULT_SCALAR)
	integral := big.NewInt(int64(x))
	integral.Mul(integral, scalar)
	return BigDec{integral: integral, scalar: scalar}
}

func FloatToBigDec(xf float64, scalar int) BigDec {
	xfBig := big.NewFloat(xf * float64(scalar))
	x := new(big.Int)
	xfBig.Int(x)
	return BigDec{integral: x, scalar: big.NewInt(int64(scalar))}
}

func (x *BigDec) ToFloat() float64 {
	// return float64(float64(x.integral) / float64(x.scalar))
	integralFloat, _ := x.integral.Float64()
	scalarFloat, _ := x.scalar.Float64()
	return integralFloat / scalarFloat
}

// TODO: change this function to return big.Int
// func (x *BigDec) GetScaledInt() int {
// 	// return int(x.integral / x.scalar)
// 	res := new(big.Int)
// 	res.Div(x.integral, x.scalar)
// 	return int(res.Int64())
// }

func FloatArrayToBigDec(array []float64, scalar int) []BigDec {
	res := make([]BigDec, len(array))
	for idx, f := range array {
		res[idx] = FloatToBigDec(f, scalar)
	}
	return res
}
