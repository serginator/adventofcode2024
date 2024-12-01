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

	result, err = Part2("input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part2 result: ", result)
}

func ReadNumbersFromFile(filename string) ([]int, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nums1 := make([]int, 0)
	nums2 := make([]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		nums := SplitLine(line)
		nums1 = append(nums1, nums[0])
		nums2 = append(nums2, nums[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return nums1, nums2, nil
}

func SplitLine(line string) []int {
	nums := make([]int, 2)
	fields := strings.Fields(line)
	if len(fields) >= 2 {
		nums[0], _ = strconv.Atoi(fields[0])
		nums[1], _ = strconv.Atoi(fields[1])
	}
	return nums
}
