package ppml

import "math/rand"

const MAX_BOOL = 2
const DEFAULT_SCALAR = 100

type BigDec struct {
	integral int64
	// We use the convention that "scalar" is the inverse of the scalar, so scalar=100 represents a 1/100th scale
	// - meaning that we divide by the scalar instead of multiplying when the representation is needed
	scalar int64
}

func CheckScalarDifference(x BigDec, y BigDec) {
	if x.scalar != y.scalar {
		panic("Different scalars not implemented")
	}
}

func Add(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	return BigDec{integral: x.integral + y.integral, scalar: x.scalar}
}

func Sub(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	return BigDec{integral: x.integral - y.integral, scalar: x.scalar}
}

func Mul(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	return BigDec{integral: (x.integral / x.scalar) * y.integral, scalar: x.scalar}
}

func Div(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	return BigDec{integral: x.integral / y.integral * x.scalar, scalar: x.scalar}
}

func Mod(x BigDec, y BigDec) BigDec {
	CheckScalarDifference(x, y)
	return BigDec{integral: (x.integral / x.scalar) % (y.integral / y.scalar) * x.scalar, scalar: x.scalar}
}

func RandBoolBigDec() BigDec {
	return BigDec{integral: int64(rand.Intn(MAX_BOOL)) * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
}

func One() BigDec {
	return BigDec{integral: 1 * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
}

func Zero() BigDec {
	return BigDec{integral: 0 * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
}

func IntToBigDecDefaultScalar(x int) BigDec {
	return BigDec{integral: int64(x) * DEFAULT_SCALAR, scalar: DEFAULT_SCALAR}
}

func FloatToBigDec(xf float32, scalar int) BigDec {
	x := int(xf * float32(scalar))
	return BigDec{int64(x), int64(scalar)}
}

func (x *BigDec) ToFloat() float32 {
	return float32(float32(x.integral) / float32(x.scalar))
}

func (x *BigDec) GetScaledInt() int {
	return int(x.integral / x.scalar)
}