# word_wrap_algorithm

Break text into lines of a maximum width using either a greedy or an optimal-fit wrapping strategy.

```go
package main

import (
	"fmt"

	wordwrap "example.com/word_wrap_algorithm"
)

func main() {
	text := "the quick brown fox jumps over the lazy dog"
	for _, line := range wordwrap.Greedy(text, 10) {
		fmt.Println(line)
	}
	// the quick
	// brown fox
	// jumps over
	// the lazy
	// dog

	for _, line := range wordwrap.Optimal(text, 10) {
		fmt.Println(line)
	}
	// the quick
	// brown fox
	// jumps over
	// the lazy
	// dog
}
```

The two exported entry points are `Greedy(text string, width int) []string` and `Optimal(text string, width int) []string`. Both also have `Options` variants, `GreedyOptions` and `OptimalOptions`, which take an `Options` struct with `Width` and `BreakLongWords` fields.

## Why this exists

Most text wrappers are greedy: fill a line, then start a new one. That is fast and predictable, but it can leave a very short final line ("a widow") when a slightly different break would have balanced the paragraph. This library offers both. `Greedy` is the simple, linear-time choice. `Optimal` uses the same dynamic-programming cost model as TeX's paragraph breaker: the cost of a line is the square of its slack (the unused width), the last line is free, and the algorithm picks the set of breaks with the lowest total cost. The trade-off is that `Optimal` is O(n²) in the number of tokens and always breaks words longer than the width, while `Greedy` is O(n) and leaves long words intact unless you ask it not to.

## The awkward edge

A word longer than the width has no good home. `Greedy` places it alone on its own line and lets it overflow, unless you set `BreakLongWords: true`, in which case it is split into chunks of exactly the width. `Optimal` always splits long words, because the cost model has no way to represent an overflowing line. If you need a long URL or hash to stay intact, use `Greedy` with the default options and accept the overflow.
