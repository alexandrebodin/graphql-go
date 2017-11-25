package lexer

import (
	"testing"
)

func TestUnhexOutOfRange(t *testing.T) {
	outOfRange := []byte{'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

	for _, b := range outOfRange {
		_, ok := unhex(b)

		if ok {
			t.Error("Expected not to be")
		}
	}
}

func TestUnhexHandlesNumber(t *testing.T) {
	numbers := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	for j := 0; j < 10; j++ {
		x, ok := unhex(numbers[j])

		if !ok || int(x) != j {
			t.Error("Expected Ok")
		}
	}
}
