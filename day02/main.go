package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	result, err := Part1("input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part1 result: ", result)

	// result, err = Part2("input")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Part2 result: ", result)
}

func ReadRowsFromFile(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rows := make([][]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		row := SplitLine(line)
		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}

func SplitLine(line string) []int {
	fields := strings.Fields(line)

	nums := make([]int, 0, len(fields))
	for _, field := range fields {
		num, err := strconv.Atoi(field)
		if err == nil {
			if num < 0 {
				log.Printf("Negative number %d found in input", num)
				return nil
			}
			nums = append(nums, num)
		} else {
			log.Printf("Error converting %s to int: %v", field, err)
			return nil
		}
	}

	return nums
}
