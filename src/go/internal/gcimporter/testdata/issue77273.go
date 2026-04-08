// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is compiled and then re-imported by TestIssue77273 to verify that
// generic methods (methods with their own type parameters) survive the
// export/import round-trip correctly.

package issue77273

// Stream is a generic type whose methods may introduce additional type
// parameters independent of the receiver type parameter T.
type Stream[T any] []T

// Map maps each element of type T to type R using f.
func (s Stream[T]) Map[R any](f func(T) R) Stream[R] {
	r := make(Stream[R], len(s))
	for i, x := range s {
		r[i] = f(x)
	}
	return r
}

// Reduce folds the stream to a single value of type R.
func (s Stream[T]) Reduce[R any](init R, f func(R, T) R) R {
	acc := init
	for _, x := range s {
		acc = f(acc, x)
	}
	return acc
}

// Filter keeps elements for which pred returns true.
func (s Stream[T]) Filter[U interface{ ~bool }](pred func(T) U) Stream[T] {
	var r Stream[T]
	for _, x := range s {
		if pred(x) {
			r = append(r, x)
		}
	}
	return r
}

// Pair holds two values of potentially different types.
type Pair[A, B any] struct {
	First  A
	Second B
}

// ZipWith combines each element with a value produced by f.
func (s Stream[T]) ZipWith[R any](f func(int, T) R) []Pair[T, R] {
	out := make([]Pair[T, R], len(s))
	for i, x := range s {
		out[i] = Pair[T, R]{x, f(i, x)}
	}
	return out
}
