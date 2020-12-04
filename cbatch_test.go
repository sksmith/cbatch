package cbatch

import (
	"bufio"
	"strings"
	"testing"
)

func TestScanMultiLines(t *testing.T) {
	const testData = `one
two
three
four
five
six
seven`

	type test struct {
		input int
		want  []string
	}

	tests := []test{
		{input: 1, want: []string{"one", "two", "three", "four", "five", "six", "seven"}},
		{input: 2, want: []string{"one\ntwo", "three\nfour", "five\nsix", "seven"}},
		{input: 0, want: []string{"one\ntwo\nthree\nfour\nfive\nsix\nseven"}},
		{input: -1, want: []string{"one\ntwo\nthree\nfour\nfive\nsix\nseven"}},
		{input: 7, want: []string{"one\ntwo\nthree\nfour\nfive\nsix\nseven"}},
		{input: 8, want: []string{"one\ntwo\nthree\nfour\nfive\nsix\nseven"}},
	}

	for _, tc := range tests {
		scanner := bufio.NewScanner(strings.NewReader(testData))
		scanner.Split(ScanMultiLines(tc.input))

		for i := 0; i < len(tc.want); i++ {
			if !scanner.Scan() {
				t.Errorf("got recordCount=[%d] want=[%d]", i, len(tc.want))
				break
			}
			if scanner.Text() != tc.want[i] {
				t.Errorf("got=[%s] want=[%s]", scanner.Text(), tc.want[i])
			}
		}

		if scanner.Scan() {
			t.Errorf("got more records than expected")
		}
	}
}

func TestScanMultiLinesEmptyData(t *testing.T) {
	const testData = ``

	type test struct {
		input int
		want  []string
	}

	tests := []test{
		{input: 1, want: []string{}},
		{input: 0, want: []string{}},
		{input: -1, want: []string{}},
		{input: 2, want: []string{}},
	}

	for _, tc := range tests {
		scanner := bufio.NewScanner(strings.NewReader(testData))
		scanner.Split(ScanMultiLines(tc.input))

		for i := 0; i < len(tc.want); i++ {
			if !scanner.Scan() {
				t.Errorf("got recordCount=[%d] want=[%d]", i, len(tc.want))
				break
			}
			if scanner.Text() != tc.want[i] {
				t.Errorf("got=[%s] want=[%s]", scanner.Text(), tc.want[i])
			}
		}

		if scanner.Scan() {
			t.Errorf("got more records than expected")
		}
	}
}
