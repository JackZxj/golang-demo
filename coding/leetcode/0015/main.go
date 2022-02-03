package main

import (
	"fmt"
	"reflect"
	"sort"
)

func main() {
	tests := []struct {
		nums []int
		want [][]int
	}{
		{[]int{}, [][]int{}},
		{[]int{0}, [][]int{}},
		{[]int{1, 1, 1, 1}, [][]int{}},
		{[]int{-1, -1, -1, -1}, [][]int{}},
		{[]int{-1, 1, 1, 1}, [][]int{}},
		{[]int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		{[]int{-1, 0, 0, 0, 1}, [][]int{{0, 0, 0}, {-1, 0, 1}}},
		{[]int{-2, 0, 1, 1, 2}, [][]int{{-2, 0, 2}, {-2, 1, 1}}},
		{[]int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{[]int{-2, -1, -1, 0, 1, 2, 2, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}, {-4, 2, 2}, {-2, 0, 2}}},
		{[]int{-6, -5, -4, -3, -2, -1, -1, 0, 1, 2, 3, 3, 4, 5}, [][]int{{-6, 1, 5}, {-6, 2, 4}, {-6, 3, 3}, {-5, 0, 5}, {-5, 1, 4}, {-5, 2, 3}, {-4, -1, 5}, {-4, 0, 4}, {-4, 1, 3}, {-3, -2, 5}, {-3, -1, 4}, {-3, 0, 3}, {-3, 1, 2}, {-2, -1, 3}, {-2, 0, 2}, {-1, -1, 2}, {-1, 0, 1}}},
	}
	for _, test := range tests {
		got := threeSum(test.nums)
		if !isEqual(got, test.want) {
			fmt.Printf("%+v,got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func isEqual(got, want [][]int) bool {
	if len(got) != len(want) {
		return false
	}
	var w [][]int
	for i := 0; i < len(want); i++ {
		t := make([]int, len(want[i]))
		copy(t, want[i])
		sort.Ints(t)
		w = append(w, t)
	}

LOOP:
	for _, nums := range got {
		sort.Ints(nums)
		for i := 0; i < len(w); i++ {
			if reflect.DeepEqual(nums, w[i]) {
				w = append(w[:i], w[i+1:]...)
				continue LOOP
			}
		}
		return false
	}
	return len(w) == 0
}

func threeSum(nums []int) [][]int {
	var result [][]int
	sort.Ints(nums)
	if len(nums) < 3 || nums[0] > 0 || nums[len(nums)-1] < 0 {
		return result
	}
	for low := 0; low < len(nums)-2; low++ {
		if nums[low] > 0 {
			break
		}
		mid, high := low+1, len(nums)-1
		for mid < high {
			if nums[low]+nums[high]+nums[mid] == 0 {
				result = append(result, []int{nums[low], nums[mid], nums[high]})
				for ; mid < high; mid++ {
					if nums[mid] != nums[mid+1] {
						mid++
						break
					}
				}
				for ; mid < high; high-- {
					if nums[high] != nums[high-1] {
						high--
						break
					}
				}
			} else if nums[low]+nums[high]+nums[mid] < 0 {
				mid++
			} else {
				high--
			}
		}
		for ; low+1 < high; low++ {
			if nums[low] != nums[low+1] {
				break
			}
		}
	}
	return result
}
