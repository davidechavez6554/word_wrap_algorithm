package wordwrap

import (
	"reflect"
	"testing"
)

func TestGreedyBasic(t *testing.T) {
	got := Greedy("the quick brown fox jumps over the lazy dog", 10)
	want := []string{
		"the quick",
		"brown fox",
		"jumps over",
		"the lazy",
		"dog",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Greedy = %#v, want %#v", got, want)
	}
}

func TestGreedyEmpty(t *testing.T) {
	if got := Greedy("", 10); got != nil {
		t.Fatalf("Greedy empty = %#v, want nil", got)
	}
	if got := Greedy("   \t\n  ", 10); got != nil {
		t.Fatalf("Greedy whitespace = %#v, want nil", got)
	}
}

func TestGreedyWidthOne(t *testing.T) {
	got := Greedy("a b c", 1)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Greedy = %#v, want %#v", got, want)
	}
}

func TestGreedyLongWordOverflow(t *testing.T) {
	// Without BreakLongWords, a word longer than width is placed alone and
	// may overflow its line.
	got := GreedyOptions("hi supercalifragilistic bye", Options{Width: 5})
	want := []string{"hi", "supercalifragilistic", "bye"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Greedy = %#v, want %#v", got, want)
	}
}

func TestGreedyLongWordBreak(t *testing.T) {
	got := GreedyOptions("hi supercalifragilistic bye", Options{Width: 5, BreakLongWords: true})
	want := []string{"hi", "super", "calif", "ragil", "istic", "bye"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Greedy = %#v, want %#v", got, want)
	}
}

func TestGreedyZeroWidth(t *testing.T) {
	// Width below 1 is clamped to 1.
	got := Greedy("ab cd", 0)
	want := []string{"ab", "cd"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Greedy = %#v, want %#v", got, want)
	}
}

func TestOptimalBasic(t *testing.T) {
	got := Optimal("the quick brown fox jumps over the lazy dog", 10)
	want := []string{
		"the quick",
		"brown fox",
		"jumps over",
		"the lazy",
		"dog",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Optimal = %#v, want %#v", got, want)
	}
}

func TestOptimalBalanced(t *testing.T) {
	// The greedy algorithm would produce "aa bb cc" / "dd" here, leaving a
	// very short last line. The optimal algorithm should balance the lines.
	got := Optimal("aa bb cc dd", 8)
	want := []string{"aa bb cc", "dd"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Optimal = %#v, want %#v", got, want)
	}
}

func TestOptimalLongWord(t *testing.T) {
	got := Optimal("hi supercalifragilistic bye", 5)
	want := []string{"hi", "super", "calif", "ragil", "istic", "bye"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Optimal = %#v, want %#v", got, want)
	}
}

func TestOptimalEmpty(t *testing.T) {
	if got := Optimal("", 10); got != nil {
		t.Fatalf("Optimal empty = %#v, want nil", got)
	}
}

func TestOptimalWidthOne(t *testing.T) {
	got := Optimal("a b c", 1)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Optimal = %#v, want %#v", got, want)
	}
}

func TestOptimalWhitespace(t *testing.T) {
	got := Optimal("  one\ttwo\nthree  ", 7)
	want := []string{"one two", "three"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Optimal = %#v, want %#v", got, want)
	}
}
