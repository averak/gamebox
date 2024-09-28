package vector

import (
	"sort"
)

type Vector[T any] []T

func New[T any](vs []T) Vector[T] {
	return vs
}

func (v Vector[T]) Raw() []T {
	return v
}

func (v Vector[T]) IsEmpty() bool {
	return len(v) == 0
}

func (v Vector[T]) Filter(eq func(T) bool) Vector[T] {
	res := make(Vector[T], 0, len(v))
	for i := range v {
		if eq(v[i]) {
			res = append(res, v[i])
		}
	}
	return res
}

func (v Vector[T]) Append(vs ...T) Vector[T] {
	return append(v, vs...)
}

func (v Vector[T]) Upsert(val T, eq func(x, y T) bool) Vector[T] {
	for i, x := range v {
		if eq(x, val) {
			(v)[i] = val
			return v
		}
	}
	return v.Append(val)
}

func (v Vector[T]) Sort(less func(x, y T) bool) Vector[T] {
	sort.Slice(v, func(i, j int) bool {
		return less((v)[i], (v)[j])
	})
	return v
}

func (v Vector[T]) Reverse() Vector[T] {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
	return v
}

func Map[S, T any](v []S, f func(S) T) []T {
	res := make([]T, len(v))
	for i := range v {
		res[i] = f(v[i])
	}
	return res
}
