package ppml

import (
	"math/big"
	"math/rand"
)

// 2^l_D
const DEFAULT_SCALAR = 1 << 13

// 2^l
const MAX_INTEGRAL = 1 << 20

const MODULUS = 2 * MAX_INTEGRAL

const SECRET_SHARE_MAX_RAND = MODULUS
const DEALER_MAX_RAND = 1 << 7

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
	modulus := big.NewInt(MODULUS)
	resIntegral.Mod(resIntegral, modulus)
	return BigDec{integral: resIntegral, scalar: x.scalar}
}

func Sub(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	resIntegral := new(big.Int)
	resIntegral.Sub(x.integral, y.integral)
	modulus := big.NewInt(MODULUS)
	resIntegral.Mod(resIntegral, modulus)
	return BigDec{integral: resIntegral, scalar: x.scalar}
}

func isNegative(x BigDec) bool {
	// fmt.Println("x.integral: ", x.integral)
	// fmt.Println("max_integral: ", MAX_INTEGRAL)
	// fmt.Println(x.integral.Cmp(big.NewInt(MAX_INTEGRAL)))
	return x.integral.Cmp(big.NewInt(MAX_INTEGRAL)) > 0
}

func Mul(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	sign := 1
	modulus := big.NewInt(MODULUS)
	xTmp := new(big.Int).Set(x.integral)
	yTmp := new(big.Int).Set(y.integral)
	// TODO: check the case x or y = 0
	if isNegative(x) {
		sign *= -1
		// fmt.Println("Max integral: ", MAX_INTEGRAL)
		// fmt.Println("before: ", x.integral)
		xTmp.Sub(modulus, x.integral)
		// fmt.Println("after: ", x.integral)
	}
	if isNegative(y) {
		sign *= -1
		yTmp.Sub(modulus, y.integral)
	}
	resIntegral := new(big.Int)
	resIntegral.Mul(xTmp, yTmp)
	resIntegral.Div(resIntegral, x.scalar)
	resIntegral.Mod(resIntegral, modulus)
	if sign < 0 {
		resIntegral.Sub(modulus, resIntegral)
	}
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

func RandForSecretShare() BigDec {
	// return BigDec{integral: int64(rand.Intn(MAX_RAND)) * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
	return BigDec{integral: big.NewInt(rand.Int63n(SECRET_SHARE_MAX_RAND)), scalar: big.NewInt(DEFAULT_SCALAR)}
}

func RandForDealer() BigDec {
	// return BigDec{integral: int64(rand.Intn(MAX_RAND)) * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
	return BigDec{integral: big.NewInt(rand.Int63n(DEALER_MAX_RAND)), scalar: big.NewInt(DEFAULT_SCALAR)}
}

// func One() BigDec {
// 	return BigDec{integral: 1 * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
// }

// func Zero() BigDec {
// 	return BigDec{integral: 0 * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
// }

func IntToBigDecDefaultScalar(x int) BigDec {
	//TODO: handle the case where x is negative
	scalar := big.NewInt(DEFAULT_SCALAR)
	integral := big.NewInt(int64(x))
	integral.Mul(integral, scalar)
	return BigDec{integral: integral, scalar: scalar}
}

func FloatToBigDec(xf float64, scalar int) BigDec {
	modulus := big.NewInt(MODULUS)
	xfBig := big.NewFloat(xf * float64(scalar))
	x := new(big.Int)
	xfBig.Int(x)
	if xf <= 0 {
		x.Add(modulus, x)
	}
	return BigDec{integral: x, scalar: big.NewInt(int64(scalar))}
}

func (x BigDec) ToFloat() float64 {
	// return float64(float64(x.integral) / float64(x.scalar))
	// fmt.Println("x.integral: ", x.integral)
	modulus := big.NewInt(MODULUS)
	integralFloat := 0.0
	if isNegative(x) {
		// fmt.Println("not else")
		// fmt.Println(x.integral)
		// fmt.Println(big.NewInt(MAX_INTEGRAL))
		tmp := new(big.Int)
		integralFloat, _ = tmp.Sub(modulus, x.integral).Float64()
		// fmt.Println("integralFloat: ", integralFloat)
		integralFloat = 0 - integralFloat
		// fmt.Println("integralFloat: ", integralFloat)
	} else {
		// fmt.Println("else")
		integralFloat, _ = x.integral.Float64()
	}
	scalarFloat, _ := x.scalar.Float64()
	// fmt.Println("x.integral: ", x.integral)
	// fmt.Println(MAX_INTEGRAL)
	// fmt.Println("integralFloat: ", integralFloat)
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
