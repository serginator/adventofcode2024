package main

import (
	"os"
	"bufio"
)

func Part2(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var nums1 []int
	var nums2 []int
	for scanner.Scan() {
		line := scanner.Text()
		nums := SplitLine(line)
		nums1 = append(nums1, nums[0])
		nums2 = append(nums2, nums[1])
	}

	sum := 0
	for i := 0; i < len(nums1); i++ {
		sum += nums1[i] * Count(nums1[i], nums2)
	}
	return sum, nil
}

func Count(num int, nums []int) int {
	count := 0
	for i := 0; i < len(nums); i++ {
		if num == nums[i] {
			count++
		}
	}
	return count
}
