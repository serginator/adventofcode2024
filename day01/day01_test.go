package main

import (
	"os"
	"testing"
)

func TestSplitLine(t *testing.T) {
	tests := []struct {
			name     string
			input    string
			expected []int
		}{
			{
				name:     "basic split",
				input:    "3   4",
				expected: []int{3, 4},
			},
			{
				name:     "different numbers",
				input:    "1   9",
				expected: []int{1, 9},
			},
			{
				name:     "same numbers",
				input:    "5   5",
				expected: []int{5, 5},
			},
			{
				name:     "multiple spaces",
				input:    "3     4",
				expected: []int{3, 4},
			},
			{
				name:     "tab separated",
				input:    "3\t4",
				expected: []int{3, 4},
			},
			{
				name:     "larger numbers",
				input:    "12  34",
				expected: []int{12, 34},
			},
		}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitLine(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitLine() returned slice of length %d, expected %d", len(result), len(tt.expected))
			}
			for i := 0; i < len(tt.expected); i++ {
				if result[i] != tt.expected[i] {
					t.Errorf("splitLine() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func TestProcess(t *testing.T) {
	tempInput := []byte(`3   4
4   3
2   5
1   3
3   9
3   3`)

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(tempInput); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := Process(tmpfile.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	expected := 11
	if result != expected {
		t.Errorf("Expected output %d, got %d", expected, result)
	}
}
