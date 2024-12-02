package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestMain(t *testing.T) {
	tempInput := []byte(`7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`)

	tmpfile, err := os.CreateTemp("", "input")
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

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	originalInput := "input"
	if _, err := os.Stat(originalInput); err == nil {
		tmpOriginal, err := os.CreateTemp("", "original_input")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpOriginal.Name())

		originalContent, err := os.ReadFile(originalInput)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(tmpOriginal.Name(), originalContent, 0644); err != nil {
			t.Fatal(err)
		}

		defer func() {
			os.Remove(originalInput)
			os.Rename(tmpOriginal.Name(), originalInput)
		}()
	}

	if err := os.Rename(tmpfile.Name(), originalInput); err != nil {
		t.Fatal(err)
	}

	main()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	expectedOutput := "Part1 result:  2\nPart2 result:  4\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestSplitLine_Errors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "non-numeric input",
			input:    "1 a 3",
			expected: nil,
		},
		{
			name:     "mixed valid and invalid numbers",
			input:    "1 2.5 3 abc 4",
			expected: nil,
		},
		{
			name:     "special characters",
			input:    "1 @#$ 2 % 3",
			expected: nil,
		},
		{
			name:     "negative numbers",
			input:    "-1 2 -3",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitLine(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitLine(%q) = %v, want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "basic split",
			input:    "3 4",
			expected: []int{3, 4},
		},
		{
			name:     "different numbers",
			input:    "1 9 6",
			expected: []int{1, 9, 6},
		},
		{
			name:     "same numbers",
			input:    "5 5",
			expected: []int{5, 5},
		},
		{
			name:     "larger numbers",
			input:    "12 34 22 81",
			expected: []int{12, 34, 22, 81},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitLine(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("SplitLine() returned slice of length %d, expected %d", len(result), len(tt.expected))
			}
			for i := 0; i < len(tt.expected); i++ {
				if result[i] != tt.expected[i] {
					t.Errorf("SplitLine() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func TestReadRowsFromFile_Errors(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "nonexistent file",
			setup: func() string {
				return "nonexistentfile.txt"
			},
			expectError: true,
			errorMsg:    "should error on nonexistent file",
		},
		{
			name: "empty file",
			setup: func() string {
				tmpfile, err := os.CreateTemp("", "empty")
				if err != nil {
					t.Fatal(err)
				}
				tmpfile.Close()
				return tmpfile.Name()
			},
			expectError: false,
			errorMsg:    "should handle empty file",
		},
		{
			name: "invalid number format",
			setup: func() string {
				tmpfile, err := os.CreateTemp("", "invalid")
				if err != nil {
					t.Fatal(err)
				}
				content := []byte("abc def\n1 2")
				if _, err := tmpfile.Write(content); err != nil {
					t.Fatal(err)
				}
				tmpfile.Close()
				return tmpfile.Name()
			},
			expectError: false,
			errorMsg:    "should handle invalid number format",
		},
		{
			name: "insufficient columns",
			setup: func() string {
				tmpfile, err := os.CreateTemp("", "insufficient")
				if err != nil {
					t.Fatal(err)
				}
				content := []byte("1\n1 2")
				if _, err := tmpfile.Write(content); err != nil {
					t.Fatal(err)
				}
				tmpfile.Close()
				return tmpfile.Name()
			},
			expectError: false,
			errorMsg:    "should handle insufficient columns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.setup()
			if filename != "nonexistentfile.txt" {
				defer os.Remove(filename)
			}

			rows, err := ReadRowsFromFile(filename)

			if tt.expectError {
				if err == nil {
					t.Error(tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if rows == nil {
					t.Error("nums1 should not be nil")
				}
				if rows == nil {
					t.Error("nums2 should not be nil")
				}
			}
		})
	}
}

func TestReadRowsFromFile_ValidContent(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		expectedRows [][]int
	}{
		{
			name:         "single line",
			content:      "1 2 4",
			expectedRows: [][]int{{1, 2, 4}},
		},
		{
			name:         "multiple lines",
			content:      "1 2\n3 4 5\n6 7 8 9",
			expectedRows: [][]int{{1, 2}, {3, 4, 5}, {6, 7, 8, 9}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatal(err)
			}
			tmpfile.Close()

			rows, err := ReadRowsFromFile(tmpfile.Name())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(rows, tt.expectedRows) {
				t.Errorf("nums1 = %v, want %v", rows, tt.expectedRows)
			}
		})
	}
}

func TestIsInOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{
			name:     "increasing sequence",
			input:    []int{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "decreasing sequence",
			input:    []int{5, 4, 3, 2, 1},
			expected: true,
		},
		{
			name:     "equal numbers",
			input:    []int{2, 2, 2, 2},
			expected: true,
		},
		{
			name:     "not in order",
			input:    []int{1, 3, 2, 4},
			expected: false,
		},
		{
			name:     "single number is in order",
			input:    []int{1},
			expected: true,
		},
		{
			name:     "two increasing numbers",
			input:    []int{1, 2},
			expected: true,
		},
		{
			name:     "two decreasing numbers",
			input:    []int{2, 1},
			expected: true,
		},
		{
			name:     "two equal numbers",
			input:    []int{2, 2},
			expected: true,
		},
		{
			name:     "mixed sequence",
			input:    []int{1, 2, 2, 1, 3},
			expected: false,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInOrder(tt.input)
			if result != tt.expected {
				t.Errorf("IsInOrder(%v) = %v, want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "positive number",
			input:    5,
			expected: 5,
		},
		{
			name:     "negative number",
			input:    -5,
			expected: 5,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Abs(tt.input)
			if result != tt.expected {
				t.Errorf("Abs(%d) = %d, want %d",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSafe(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{
			name:     "increasing sequence",
			input:    []int{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "decreasing sequence",
			input:    []int{5, 4, 3, 2, 1},
			expected: true,
		},
		{
			name:     "equal numbers is not safe",
			input:    []int{2, 2, 2, 2},
			expected: false,
		},
		{
			name:     "not in order",
			input:    []int{1, 3, 2, 4},
			expected: false,
		},
		{
			name:     "single number is in order",
			input:    []int{1},
			expected: true,
		},
		{
			name:     "two increasing numbers",
			input:    []int{1, 2},
			expected: true,
		},
		{
			name:     "two decreasing numbers",
			input:    []int{2, 1},
			expected: true,
		},
		{
			name:     "two equal numbers is not safe",
			input:    []int{2, 2},
			expected: false,
		},
		{
			name:     "mixed sequence",
			input:    []int{1, 2, 2, 1, 3},
			expected: false,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: true,
		},
		{
			name:     "large difference",
			input:    []int{1, 5, 2, 4},
			expected: false,
		},
		{
			name:     "small difference",
			input:    []int{1, 2, 3, 4},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSafe(tt.input)
			if result != tt.expected {
				t.Errorf("IsSafe(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPart1Error(t *testing.T) {
	_, err := Part1("nonexistentfile")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestPart1(t *testing.T) {
	tempInput := []byte(`7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`)

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

	result, err := Part1(tmpfile.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	expected := 2
	if result != expected {
		t.Errorf("Expected output %d, got %d", expected, result)
	}
}

func TestCanBeMadeSafe(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{
			name:     "already safe sequence",
			input:    []int{7, 6, 4, 2, 1},
			expected: true,
		},
		{
			name:     "cannot be made safe",
			input:    []int{1, 2, 7, 8, 9},
			expected: false,
		},
		{
			name:     "another cannot be made safe",
			input:    []int{9, 7, 6, 2, 1},
			expected: false,
		},
		{
			name:     "can be made safe by removing one element",
			input:    []int{1, 3, 2, 4, 5},
			expected: true,
		},
		{
			name:     "can be made safe by removing duplicate",
			input:    []int{8, 6, 4, 4, 1},
			expected: true,
		},
		{
			name:     "already safe increasing",
			input:    []int{1, 3, 6, 7, 9},
			expected: true,
		},
		{
			name:     "single element",
			input:    []int{1},
			expected: true,
		},
		{
			name:     "two elements safe",
			input:    []int{2, 1},
			expected: true,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: true,
		},
		{
			name:     "three elements can be made safe",
			input:    []int{1, 5, 2},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanBeMadeSafe(tt.input)
			if result != tt.expected {
				t.Errorf("CanBeMadeSafe(%v) = %v, want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestPart2Error(t *testing.T) {
	_, err := Part2("nonexistentfile")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	tmpfile, err := os.CreateTemp("", "empty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	result, err := Part2(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unexpected error for empty file: %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0 for empty file, got %d", result)
	}

	invalidInput := []byte("invalid content\n")
	tmpfileInvalid, err := os.CreateTemp("", "invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfileInvalid.Name())

	if _, err := tmpfileInvalid.Write(invalidInput); err != nil {
		t.Fatal(err)
	}
	if err := tmpfileInvalid.Close(); err != nil {
		t.Fatal(err)
	}

	result, err = Part2(tmpfileInvalid.Name())
	if err != nil {
		t.Fatalf("Unexpected error for invalid content: %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0 for invalid content, got %d", result)
	}
}

func TestPart2(t *testing.T) {
	tempInput := []byte(`7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`)

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

	result, err := Part2(tmpfile.Name())
	if err != nil {
		t.Fatalf("Part2 failed: %v", err)
	}

	expected := 4
	if result != expected {
		t.Errorf("Expected output %d, got %d", expected, result)
	}
}
