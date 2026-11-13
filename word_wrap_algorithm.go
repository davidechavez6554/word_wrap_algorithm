// Package wordwrap provides greedy and optimal-fit line wrapping for text.
//
// The package name is wordwrap, matching the module name in go.mod. The
// exported API is intentionally small: two wrapping functions and a small
// options struct.
package wordwrap

import (
	"math"
	"strings"
)

// Options controls how text is wrapped.
type Options struct {
	// Width is the maximum number of characters allowed on a line. It must be
	// at least 1; values below 1 are treated as 1.
	Width int
	// BreakLongWords, when true, splits words longer than Width across lines
	// instead of letting them overflow. The greedy algorithm only breaks words
	// when this is true. The optimal algorithm always breaks long words because
	// otherwise a single long word would make the cost model meaningless.
	BreakLongWords bool
}

// Greedy wraps text using the classic greedy line-breaking algorithm: it
// fills each line with as many words as fit in Width, then starts a new line.
//
// Words are separated by runs of ASCII whitespace (space, tab, newline, etc.)
// as defined by strings.Fields. The returned slice never contains an empty
// line unless the input contains no words at all, in which case it is empty.
func Greedy(text string, width int) []string {
	return GreedyOptions(text, Options{Width: width})
}

// GreedyOptions is Greedy with explicit options.
func GreedyOptions(text string, opts Options) []string {
	w := opts.Width
	if w < 1 {
		w = 1
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var current strings.Builder
	currentLen := 0

	flush := func() {
		if currentLen > 0 {
			lines = append(lines, current.String())
			current.Reset()
			currentLen = 0
		}
	}

	for _, word := range words {
		// A word that does not fit on an empty line can only be placed by
		// breaking it, if the caller allows that.
		if len(word) > w {
			if !opts.BreakLongWords {
				// Overflow: put the word on its own line. This is the
				// least surprising behaviour for callers who did not ask
				// for word breaking.
				flush()
				lines = append(lines, word)
				continue
			}
			// Break the word into chunks of at most w characters.
			flush()
			for len(word) > w {
				lines = append(lines, word[:w])
				word = word[w:]
			}
			if len(word) > 0 {
				current.WriteString(word)
				currentLen = len(word)
			}
			continue
		}

		if currentLen == 0 {
			current.WriteString(word)
			currentLen = len(word)
			continue
		}
		if currentLen+1+len(word) <= w {
			current.WriteByte(' ')
			current.WriteString(word)
			currentLen += 1 + len(word)
			continue
		}
		flush()
		current.WriteString(word)
		currentLen = len(word)
	}
	flush()
	return lines
}

// Optimal wraps text using a dynamic programming algorithm that minimises the
// sum of squared slack on every line except the last. This is the same cost
// model used by TeX's paragraph-breaking algorithm: a line that is almost
// full costs little, a line that is nearly empty costs a lot, and the final
// line is free so the paragraph can end naturally.
//
// Words are separated by runs of ASCII whitespace. Long words are broken
// across lines because the cost model requires every word to fit on a line.
// The returned slice never contains an empty line unless the input contains
// no words at all, in which case it is empty.
func Optimal(text string, width int) []string {
	return OptimalOptions(text, Options{Width: width})
}

// OptimalOptions is Optimal with explicit options. The BreakLongWords field
// is ignored: optimal wrapping always breaks long words, because otherwise a
// word longer than Width would have no valid placement in the cost model.
func OptimalOptions(text string, opts Options) []string {
	w := opts.Width
	if w < 1 {
		w = 1
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	// Break any word longer than w into pieces. Each piece is a separate
	// token; the pieces are joined without spaces when they came from the
	// same original word, but since we only join tokens with spaces in the
	// DP below, we track whether a token is a continuation.
	type token struct {
		text         string
		continuation bool // true if this token continues the previous token
	}
	var tokens []token
	for _, word := range words {
		if len(word) <= w {
			tokens = append(tokens, token{text: word})
			continue
		}
		first := true
		for len(word) > w {
			tokens = append(tokens, token{text: word[:w], continuation: !first})
			word = word[w:]
			first = false
		}
		if len(word) > 0 {
			tokens = append(tokens, token{text: word, continuation: !first})
		}
	}

	n := len(tokens)
	// cost[i] is the minimum total cost of wrapping tokens[i:].
	// next[i] is the index of the first token on the line that starts at i.
	cost := make([]float64, n+1)
	next := make([]int, n+1)
	for i := range cost {
		cost[i] = math.Inf(1)
	}
	cost[n] = 0

	for i := n - 1; i >= 0; i-- {
		lineLen := 0
		for j := i; j < n; j++ {
			// Compute the length of the line containing tokens i..j.
			// Continuation tokens are appended without a space.
			if j == i {
				lineLen = len(tokens[j].text)
			} else if tokens[j].continuation {
				lineLen += len(tokens[j].text)
			} else {
				lineLen += 1 + len(tokens[j].text)
			}
			if lineLen > w {
				break
			}
			// The last line is free, so we do not add its slack to the cost.
			var lineCost float64
			if j == n-1 {
				lineCost = 0
			} else {
				slack := w - lineLen
				lineCost = float64(slack * slack)
			}
			total := lineCost + cost[j+1]
			if total < cost[i] {
				cost[i] = total
				next[i] = j + 1
			}
		}
	}

	var lines []string
	for i := 0; i < n; {
		j := next[i]
		if j <= i {
			// Should not happen because every token fits on a line by
			// itself, but guard against an infinite loop.
			j = i + 1
		}
		var b strings.Builder
		for k := i; k < j; k++ {
			if k > i && !tokens[k].continuation {
				b.WriteByte(' ')
			}
			b.WriteString(tokens[k].text)
		}
		lines = append(lines, b.String())
		i = j
	}
	return lines
}
