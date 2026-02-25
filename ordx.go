package ordx

import "fmt"

// AsLess converts a three-way comparison function into a less function.
//
// The provided cmp function must return:
//   - -1 if a < b
//   - 0 if a == b
//   - +1 if a > b
//
// The returned function reports whether a is strictly less than b.
func AsLess[T any](cmp func(T, T) int) func(T, T) bool {
	return func(a, b T) bool {
		return cmp(a, b) == -1
	}
}

// AsLess converts a three-way comparison function into a less function.
//
// The provided cmp function must return:
//   - -1 if a < b
//   - 0 if a == b
//   - +1 if a > b
//
// The returned function reports whether a is strictly less than b.
func AsCmp[T any](less func(T, T) bool) func(T, T) int {
	return func(a, b T) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return +1
		default:
			return 0
		}
	}
}

// RankCmp returns a comparison function based on an explicit ordering.
//
// The order slice defines the ranking from lowest to highest.
// All values must be unique; duplicates cause a panic.
// Comparing a value not present in the order also causes a panic.
//
// Example:
//
//	cmp := RankCmp([]string{"low", "med", "high"})
//	cmp("low", "high")   // -1
//	cmp("high", "low")   // 1
//	cmp("med", "med")    // 0
//
// See a safer alternative [RankCmpSafe].
func RankCmp[S ~[]E, E comparable](order S) func(a, b E) int {
	m := make(map[E]int, len(order))
	for i, v := range order {
		if _, ok := m[v]; ok {
			panic(fmt.Sprintf("ordx: duplicated element in order: '%+v'", v))
		}
		m[v] = i
	}

	return func(a, b E) int {
		ai, aok := m[a]
		bi, bok := m[b]
		switch {
		case !aok:
			panic(fmt.Sprintf("ordx: unknown value: '%+v'", a))
		case !bok:
			panic(fmt.Sprintf("ordx: unknown value: '%+v'", b))
		case ai == bi:
			return 0
		case ai < bi:
			return -1
		default:
			return 1
		}
	}
}

// RankCmpSafe returns a comparison function based on an explicit ordering.
//
// The order slice defines the ranking from lowest to highest.
// All values must be unique; duplicates return an error.
// Comparing a value not present in the order also causes a panic.
//
// This is a safer alternative to [RankCmp] that validates the input slice
// instead of panicking on duplicates during construction.
func RankCmpSafe[S ~[]E, E comparable](order S) (func(a, b E) int, error) {
	m := make(map[E]int, len(order))
	for i, v := range order {
		if _, ok := m[v]; ok {
			return nil, fmt.Errorf("ordx: duplicated element in order: '%+v'", v)
		}
		m[v] = i
	}

	return func(a, b E) int {
		ai, aok := m[a]
		bi, bok := m[b]
		switch {
		case !aok:
			panic(fmt.Sprintf("ordx: unknown value: '%+v'", a))
		case !bok:
			panic(fmt.Sprintf("ordx: unknown value: '%+v'", b))
		case ai < bi:
			return -1
		case ai > bi:
			return 1
		default:
			return 0
		}
	}, nil
}

// ChainCmp combines multiple comparison functions.
//
// Each comparator is applied in order until one returns a non-zero result.
// If all comparators report equality, the result is 0.
func ChainCmp[T any](cmps ...func(T, T) int) func(a, b T) int {
	return func(a, b T) int {
		for _, cmp := range cmps {
			if c := cmp(a, b); c != 0 {
				return c
			}
		}
		return 0
	}
}
