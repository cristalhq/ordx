package ordx_test

import (
	"strings"
	"testing"

	"github.com/cristalhq/ordx"
)

func TestAsLess_Int(t *testing.T) {
	cmp := func(a, b int) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	}

	less := ordx.AsLess(cmp)

	tests := []struct {
		name string
		a, b int
		want bool
	}{
		{"a < b", 1, 2, true},
		{"a > b", 2, 1, false},
		{"a == b", 2, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := less(tt.a, tt.b); got != tt.want {
				t.Fatalf("less(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAsLess_String(t *testing.T) {
	cmp := func(a, b string) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	}

	less := ordx.AsLess(cmp)

	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{`"a" < "b"`, "a", "b", true},
		{`"b" < "a"`, "b", "a", false},
		{`"a" == "a"`, "a", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := less(tt.a, tt.b); got != tt.want {
				t.Fatalf("less(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAsCmp_Int(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	cmp := ordx.AsCmp(less)

	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"a < b", 1, 2, -1},
		{"a > b", 2, 1, 1},
		{"a == b", 2, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmp(tt.a, tt.b); got != tt.want {
				t.Fatalf("cmp(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAsCmp_String(t *testing.T) {
	less := func(a, b string) bool { return a < b }
	cmp := ordx.AsCmp(less)

	tests := []struct {
		name string
		a, b string
		want int
	}{
		{`"a" < "b"`, "a", "b", -1},
		{`"b" > "a"`, "b", "a", 1},
		{`"a" == "a"`, "a", "a", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmp(tt.a, tt.b); got != tt.want {
				t.Fatalf("cmp(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestRankCmp(t *testing.T) {
	t.Run("BasicOrdering", func(t *testing.T) {
		order := []string{"a", "b", "c"}
		cmp := ordx.RankCmp(order)

		tests := []struct {
			a, b     string
			expected int
		}{
			{"a", "a", 0},
			{"a", "b", -1},
			{"b", "a", 1},
			{"b", "c", -1},
			{"c", "b", 1},
		}

		for _, tt := range tests {
			got := cmp(tt.a, tt.b)
			if got != tt.expected {
				t.Fatalf("cmp(%q,%q) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		}
	})

	t.Run("Ints", func(t *testing.T) {
		order := []int{10, 20, 30}
		cmp := ordx.RankCmp(order)

		if cmp(10, 30) != -1 {
			t.Fatalf("expected 10 < 30")
		}
		if cmp(20, 20) != 0 {
			t.Fatalf("expected 20 == 20")
		}
		if cmp(30, 10) != 1 {
			t.Fatalf("expected 30 > 10")
		}
	})

	t.Run("DuplicatePanics", func(t *testing.T) {
		order := []string{"x", "y", "x"}

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic on duplicate order element")
			}
		}()

		_ = ordx.RankCmp(order)
	})

	t.Run("UnknownElementPanics_A", func(t *testing.T) {
		order := []string{"u", "v", "w"}
		cmp := ordx.RankCmp(order)

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic on unknown element (a)")
			}
		}()

		cmp("zzz", "u") // a is unknown
	})

	t.Run("UnknownElementPanics_B", func(t *testing.T) {
		order := []string{"u", "v", "w"}
		cmp := ordx.RankCmp(order)

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic on unknown element (b)")
			}
		}()

		cmp("u", "zzz") // b is unknown
	})

	t.Run("Stability", func(t *testing.T) {
		order := []string{"k", "l", "m"}
		cmp := ordx.RankCmp(order)

		if cmp("k", "m") != -1 {
			t.Fatalf("expected k < m")
		}
		if cmp("m", "l") != 1 {
			t.Fatalf("expected m > l")
		}
	})
}

func TestRankCmpSafe(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		cmp, err := ordx.RankCmpSafe([]string{"low", "med", "high"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tests := []struct {
			a, b string
			want int
		}{
			{"low", "low", 0},
			{"low", "med", -1},
			{"med", "low", 1},
			{"med", "high", -1},
			{"high", "med", 1},
		}

		for _, tt := range tests {
			got := cmp(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("cmp(%q, %q) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		}
	})

	t.Run("Duplicate", func(t *testing.T) {
		_, err := ordx.RankCmpSafe([]string{"a", "b", "a"})
		if err == nil {
			t.Fatalf("expected error for duplicate elements, got nil")
		}
	})

	t.Run("UnknownPanics_A", func(t *testing.T) {
		cmp, err := ordx.RankCmpSafe([]string{"a", "b"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic for unknown value")
			}
		}()

		_ = cmp("a", "c")
	})

	t.Run("UnknownPanics_B", func(t *testing.T) {
		cmp, err := ordx.RankCmpSafe([]string{"a", "b"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic for unknown value")
			}
		}()

		_ = cmp("c", "a")
	})
}

func TestChainCmp_Lexicographic(t *testing.T) {
	type User struct {
		First string
		Last  string
	}

	cmp := ordx.ChainCmp(
		func(a, b User) int { return strings.Compare(a.Last, b.Last) },
		func(a, b User) int { return strings.Compare(a.First, b.First) },
	)

	tests := []struct {
		a, b User
		want int
	}{
		{User{"John", "Doe"}, User{"John", "Doe"}, 0},
		{User{"Alice", "Doe"}, User{"Bob", "Doe"}, -1}, // first name tiebreak
		{User{"Bob", "Doe"}, User{"Alice", "Doe"}, 1},
		{User{"John", "Adams"}, User{"John", "Doe"}, -1}, // last name primary
		{User{"John", "Doe"}, User{"John", "Adams"}, 1},
	}

	for _, tt := range tests {
		got := cmp(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("cmp(%v, %v) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func BenchmarkAsLess_Compare(b *testing.B) {
	less := ordx.AsLess(intCmp)
	a, c := 10, 20

	for b.Loop() {
		if !less(a, c) {
			b.Fatal("not less")
		}
	}
}

func BenchmarkAsCmp_Compare(b *testing.B) {
	cmp := ordx.AsCmp(intLess)
	a, c := 10, 20

	for b.Loop() {
		if v := cmp(a, c); v != -1 {
			b.Fatalf("have: %v", v)
		}
	}
}

func BenchmarkRankCmp_Compare(b *testing.B) {
	order := []int{1, 2, 3, 4, 5}
	cmp := ordx.RankCmp(order)

	a, c := 2, 4
	for b.Loop() {
		if v := cmp(a, c); v != -1 {
			b.Fatalf("have: %v", v)
		}
	}
}

func intCmp(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func intLess(a, b int) bool {
	return a < b
}
