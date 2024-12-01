package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"
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


func SplitLine(line string) []int {
	nums := make([]int, 2)
	fields := strings.Fields(line)
	if len(fields) >= 2 {
		nums[0], _ = strconv.Atoi(fields[0])
		nums[1], _ = strconv.Atoi(fields[1])
	}
	return nums
}
