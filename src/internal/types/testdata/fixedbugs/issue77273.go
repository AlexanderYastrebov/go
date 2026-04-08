// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tests for generic methods (go.dev/issue/77273).
// A method may declare its own type parameters, allowing
// mapping from the receiver's type T to a different type R.

package p

// Stream is a generic type whose methods may introduce additional type parameters.
type Stream[T any] []T

// Map maps each element of type T to type R using f.
// This is the motivating example: without generic methods, Map must be a
// free function because the result type R differs from the receiver type T.
func (s Stream[T]) Map[R any](f func(T) R) Stream[R] {
	r := make(Stream[R], len(s))
	for i, x := range s {
		r[i] = f(x)
	}
	return r
}

// Filter keeps elements satisfying pred.
func (s Stream[T]) Filter[U interface{ ~bool }](pred func(T) U) Stream[T] {
	var r Stream[T]
	for _, x := range s {
		if pred(x) {
			r = append(r, x)
		}
	}
	return r
}

// Reduce folds the stream to a single value of potentially different type R.
func (s Stream[T]) Reduce[R any](init R, f func(R, T) R) R {
	acc := init
	for _, x := range s {
		acc = f(acc, x)
	}
	return acc
}

func _() {
	s := Stream[int]{1, 2, 3}

	// Map with explicit type argument: int -> string
	r1 := s.Map[string](func(x int) string { return "x" })
	var _ Stream[string] = r1

	// Map with type inference: int -> float64
	r2 := s.Map(func(x int) float64 { return float64(x) * 1.5 })
	var _ Stream[float64] = r2

	// Chaining: Stream[int] -> Stream[string] -> Stream[int]
	r3 := s.Map[string](func(x int) string { return "y" }).Map[int](func(s string) int { return len(s) })
	var _ Stream[int] = r3

	// Reduce: int -> int
	sum := s.Reduce[int](0, func(acc, x int) int { return acc + x })
	var _ int = sum

	// Reduce with type inference on R: int -> string
	_ = s.Reduce("", func(acc string, x int) string { return acc })
}

// Pair holds two values of potentially different types.
type Pair[A, B any] struct {
	First  A
	Second B
}

// ZipWith combines each element with a value produced by f, returning a slice of Pairs.
func (s Stream[T]) ZipWith[R any](f func(int, T) R) []Pair[T, R] {
	out := make([]Pair[T, R], len(s))
	for i, x := range s {
		out[i] = Pair[T, R]{x, f(i, x)}
	}
	return out
}

func _() {
	s := Stream[int]{10, 20, 30}
	zipped := s.ZipWith[string](func(i int, x int) string { return "v" })
	var _ []Pair[int, string] = zipped
}
