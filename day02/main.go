package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	nums := make([]int, 0)
	for i := 0; i < len(line); i++ {
		if line[i] == ' ' {
			continue
		}
		num := 0
		for i < len(line) && line[i] != ' ' {
			num = num*10 + int(line[i]-'0')
			i++
		}
		nums = append(nums, num)
	}
	return nums
}
