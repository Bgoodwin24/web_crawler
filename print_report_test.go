package main

import (
	"bytes"
	"os"
	"testing"
)

func TestPrintReport(t *testing.T) {
	tests := []struct {
		name     string
		pages    map[string]int
		baseURL  string
		expected string
	}{
		{
			name: "Single page",
			pages: map[string]int{
				"https://example.com/page1": 5,
			},
			baseURL:  "https://example.com",
			expected: "=============================\n  REPORT for https://example.com\n=============================\nFound 5 internal links to https://example.com/page1\n",
		},
		{
			name: "Multiple pages, sorted by count",
			pages: map[string]int{
				"https://example.com/page2": 3,
				"https://example.com/page1": 5,
			},
			baseURL:  "https://example.com",
			expected: "=============================\n  REPORT for https://example.com\n=============================\nFound 5 internal links to https://example.com/page1\nFound 3 internal links to https://example.com/page2\n",
		},
		{
			name: "Tie in counts, sorted alphabetically",
			pages: map[string]int{
				"https://example.com/pageB": 3,
				"https://example.com/pageA": 3,
			},
			baseURL:  "https://example.com",
			expected: "=============================\n  REPORT for https://example.com\n=============================\nFound 3 internal links to https://example.com/pageA\nFound 3 internal links to https://example.com/pageB\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			originalStdout := os.Stdout
			os.Stdout = w

			printReport(tc.pages, tc.baseURL)

			w.Close()
			var buf bytes.Buffer
			buf.ReadFrom(r)
			os.Stdout = originalStdout

			if buf.String() != tc.expected {
				t.Errorf("FAIL (%s): expected\n%v\ngot\n%v", tc.name, tc.expected, buf.String())
			}
		})
	}
}
