package main

import (
	"sort"
)

func Part1(filename string) (int, error) {
	nums1, nums2, err := ReadNumbersFromFile(filename)
	if err != nil {
		return 0, err
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
