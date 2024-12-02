package main

func Part2(filename string) (int, error) {
	rows, err := ReadRowsFromFile(filename)
	if err != nil {
		return 0, err
	}

	safe := 0
	for i := 0; i < len(rows); i++ {
		if IsSafe(rows[i]) || CanBeMadeSafe(rows[i]) {
			safe++
		}
	}
	return safe, nil
}

func CanBeMadeSafe(row []int) bool {
	if len(row) <= 1 {
		return true
	}

	if IsSafe(row) {
		return true
	}

	for i := 0; i < len(row); i++ {
		newRow := make([]int, 0, len(row)-1)
		newRow = append(newRow, row[:i]...)
		newRow = append(newRow, row[i+1:]...)

		if IsSafe(newRow) {
			return true
		}
	}
	return false
}
