package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestMain(t *testing.T) {
	// Create a temporary test input file
	tempInput := []byte(`3   4
4   3
2   5
1   3
3   9
3   3`)

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

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Temporarily rename the test file to "input"
	originalInput := "input"
	if _, err := os.Stat(originalInput); err == nil {
		// If input file exists, save it
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

		// Restore original input file after test
		defer func() {
			os.Remove(originalInput)
			os.Rename(tmpOriginal.Name(), originalInput)
		}()
	}

	// Copy our test file to "input"
	if err := os.Rename(tmpfile.Name(), originalInput); err != nil {
		t.Fatal(err)
	}

	// Run main
	main()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	// Check expected output
	expectedOutput := "Part1 result:  11\nPart2 result:  31\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
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
			input:    "3	 4",
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

func TestReadNumbersFromFile_Errors(t *testing.T) {
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

			nums1, nums2, err := ReadNumbersFromFile(filename)

			if tt.expectError {
				if err == nil {
					t.Error(tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if nums1 == nil {
					t.Error("nums1 should not be nil")
				}
				if nums2 == nil {
					t.Error("nums2 should not be nil")
				}
			}
		})
	}
}

func TestReadNumbersFromFile_ValidContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedNums1 []int
		expectedNums2 []int
	}{
		{
			name:          "single line",
			content:       "1 2",
			expectedNums1: []int{1},
			expectedNums2: []int{2},
		},
		{
			name:          "multiple lines",
			content:       "1 2\n3 4\n5 6",
			expectedNums1: []int{1, 3, 5},
			expectedNums2: []int{2, 4, 6},
		},
		{
			name:          "extra whitespace",
			content:       "1	2\n3  4",
			expectedNums1: []int{1, 3},
			expectedNums2: []int{2, 4},
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

			nums1, nums2, err := ReadNumbersFromFile(tmpfile.Name())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(nums1, tt.expectedNums1) {
				t.Errorf("nums1 = %v, want %v", nums1, tt.expectedNums1)
			}
			if !reflect.DeepEqual(nums2, tt.expectedNums2) {
				t.Errorf("nums2 = %v, want %v", nums2, tt.expectedNums2)
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

	result, err := Part1(tmpfile.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	expected := 11
	if result != expected {
		t.Errorf("Expected output %d, got %d", expected, result)
	}
}

func TestPart2Error(t *testing.T) {
	_, err := Part2("nonexistentfile")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestPart2(t *testing.T) {
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

	result, err := Part2(tmpfile.Name())
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	expected := 31
	if result != expected {
		t.Errorf("Expected output %d, got %d", expected, result)
	}
}
