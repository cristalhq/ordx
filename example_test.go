package ordx_test

import (
	"fmt"
	"sort"

	"github.com/cristalhq/ordx"
)

func ExampleAsLess() {
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

	values := []int{3, 1, 2}
	sort.Slice(values, func(i, j int) bool {
		return less(values[i], values[j])
	})

	fmt.Println(values)

	// Output: [1 2 3]
}

func ExampleAsCmp() {
	less := func(a, b int) bool { return a < b }
	cmp := ordx.AsCmp(less)

	fmt.Println(cmp(1, 2))
	fmt.Println(cmp(2, 1))
	fmt.Println(cmp(2, 2))

	// Output:
	// -1
	// 1
	// 0
}

func ExampleRankCmp() {
	cmp := ordx.RankCmp([]string{"low", "medium", "high"})

	fmt.Println(cmp("low", "high"))
	fmt.Println(cmp("high", "low"))
	fmt.Println(cmp("medium", "medium"))

	// Output:
	// -1
	// 1
	// 0
}
