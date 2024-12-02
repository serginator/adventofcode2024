package main

func Part1(filename string) (int, error) {
	rows, err := ReadRowsFromFile(filename)
	if err != nil {
		return 0, err
	}

	safe := 0
	for i := 0; i < len(rows); i++ {
		if IsSafe(rows[i]) {
			safe++
		}
	}
	return safe, nil
}

func IsInOrder(row []int) bool {
	increasing := true
	decreasing := true
	for i := 1; i < len(row); i++ {
		if row[i] > row[i-1] {
			decreasing = false
		}
		if row[i] < row[i-1] {
			increasing = false
		}
	}
	if increasing || decreasing {
		return true
	}
	return false
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func IsSafe(row []int) bool {
	if !IsInOrder(row) {
		return false
	}
	for i := 1; i < len(row); i++ {
		if (row[i] == row[i-1]) || (Abs(row[i]-row[i-1]) > 3) {
			return false
		}
	}
	return true
}
