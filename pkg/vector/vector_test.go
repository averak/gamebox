package vector

import (
	"reflect"
	"testing"
)

func TestVector_IsEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		v    Vector[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "empty",
			v:    Vector[int]{},
			want: true,
		},
		{
			name: "not empty",
			v: Vector[int]{
				1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Filter(t *testing.T) {
	type args[T any] struct {
		eq func(T) bool
	}
	type testCase[T any] struct {
		name string
		v    Vector[T]
		args args[T]
		want Vector[T]
	}
	tests := []testCase[int]{
		{
			v: Vector[int]{1, 2, 3, 4, 5},
			args: args[int]{
				eq: func(i int) bool {
					return i%2 == 0
				},
			},
			want: Vector[int]{2, 4},
		},
		{
			v: Vector[int]{1, 2, 3, 4, 5},
			args: args[int]{
				eq: func(i int) bool {
					return i%2 != 0
				},
			},
			want: Vector[int]{1, 3, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Filter(tt.args.eq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Sort(t *testing.T) {
	type args[T any] struct {
		less func(x, y T) bool
	}
	type testCase[T any] struct {
		name string
		v    Vector[T]
		args args[T]
		want Vector[T]
	}
	tests := []testCase[int]{
		{
			v: Vector[int]{3, 2, 1},
			args: args[int]{
				less: func(x, y int) bool {
					return x < y
				},
			},
			want: Vector[int]{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Sort(tt.args.less); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Reverse(t *testing.T) {
	type testCase[T any] struct {
		name string
		v    Vector[T]
		want Vector[T]
	}
	tests := []testCase[int]{
		{
			v:    Vector[int]{1, 2, 3},
			want: Vector[int]{3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Append(t *testing.T) {
	type args[T any] struct {
		vs []T
	}
	type testCase[T any] struct {
		name string
		v    Vector[T]
		args args[T]
		want Vector[T]
	}
	tests := []testCase[int]{
		{
			v: Vector[int]{1, 2, 3},
			args: args[int]{
				vs: []int{4, 5, 6},
			},
			want: Vector[int]{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Append(tt.args.vs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Upsert(t *testing.T) {
	type s struct {
		ID    string
		Value int
	}

	type args[T any] struct {
		val T
		eq  func(x, y T) bool
	}
	type testCase[T any] struct {
		name string
		v    Vector[T]
		args args[T]
		want Vector[T]
	}
	tests := []testCase[s]{
		{
			name: "存在しない場合 => insert",
			v: Vector[s]{
				{
					ID:    "1",
					Value: 1,
				},
				{
					ID:    "2",
					Value: 2,
				},
			},
			args: args[s]{
				val: s{
					ID:    "3",
					Value: 3,
				},
				eq: func(x, y s) bool {
					return x.ID == y.ID
				},
			},
			want: Vector[s]{
				{
					ID:    "1",
					Value: 1,
				},
				{
					ID:    "2",
					Value: 2,
				},
				{
					ID:    "3",
					Value: 3,
				},
			},
		},
		{
			name: "存在する場合 => update",
			v: Vector[s]{
				{
					ID:    "1",
					Value: 1,
				},
				{
					ID:    "2",
					Value: 2,
				},
			},
			args: args[s]{
				val: s{
					ID:    "1",
					Value: 999,
				},
				eq: func(x, y s) bool {
					return x.ID == y.ID
				},
			},
			want: Vector[s]{
				{
					ID:    "1",
					Value: 999,
				},
				{
					ID:    "2",
					Value: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Upsert(tt.args.val, tt.args.eq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Upsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type args[S any, T any] struct {
		v []S
		f func(S) T
	}
	type testCase[S any, T any] struct {
		name string
		args args[S, T]
		want []T
	}
	tests := []testCase[int, int]{
		{
			args: args[int, int]{
				v: []int{1, 2, 3},
				f: func(i int) int {
					return i * 2
				},
			},
			want: []int{2, 4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.v, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
