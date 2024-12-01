package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Process(filename string) (int, error) {
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
		nums := splitLine(line)
		nums1 = append(nums1, nums[0])
		nums2 = append(nums2, nums[1])
	}

	sort.Ints(nums1)
	sort.Ints(nums2)

	sum := 0
	for i := 0; i < len(nums1); i++ {
		if nums1[i] >= nums2[i] {
			sum += nums1[i] - nums2[i]
		} else {
			sum += nums2[i] - nums1[i]
		}
	}
	return sum, nil
}

func splitLine(line string) []int {
	nums := make([]int, 2)
	fields := strings.Fields(line)
	if len(fields) >= 2 {
		nums[0], _ = strconv.Atoi(fields[0])
		nums[1], _ = strconv.Atoi(fields[1])
	}
	return nums
}
