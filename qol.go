package btils

func None[T any]() T {
	return *new(T)
}

func If[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}
