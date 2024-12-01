package main

import (
	"sort"
)

func Part2(filename string) (int, error) {
	nums1, nums2, err := ReadNumbersFromFile(filename)
	if err != nil {
		return 0, err
	}

	sort.Ints(nums1)
	sort.Ints(nums2)

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
