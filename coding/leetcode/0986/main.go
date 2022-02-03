package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		first  [][]int
		second [][]int
		want   [][]int
	}{
		{[][]int{{0, 2}, {5, 10}, {13, 23}, {24, 25}}, [][]int{{1, 5}, {8, 12}, {15, 24}, {25, 26}}, [][]int{{1, 2}, {5, 5}, {8, 10}, {15, 23}, {24, 24}, {25, 25}}},
		{[][]int{{1, 3}, {5, 9}}, [][]int{}, nil},
		{[][]int{}, [][]int{{4, 8}, {10, 12}}, nil},
		{[][]int{{1, 7}}, [][]int{{3, 10}}, [][]int{{3, 7}}},
		{[][]int{{3, 7}}, [][]int{{1, 10}}, [][]int{{3, 7}}},
		{[][]int{{1, 3}}, [][]int{{7, 10}}, nil},
		{[][]int{{7, 10}}, [][]int{{1, 3}}, nil},
		{[][]int{{7, 10}}, [][]int{{1, 3}, {11, 12}}, nil},
		{[][]int{{1, 3}}, [][]int{{3, 10}}, [][]int{{3, 3}}},
		{[][]int{{1, 3}, {10, 11}}, [][]int{{3, 10}}, [][]int{{3, 3}, {10, 10}}},
	}
	for _, test := range tests {
		got := intervalIntersection(test.first, test.second)
		if !reflect.DeepEqual(got, test.want) {
			fmt.Printf("%+v,got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func intervalIntersection(firstList [][]int, secondList [][]int) [][]int {
	var result [][]int
	for index1, index2 := 0, 0; index1 < len(firstList) && index2 < len(secondList); {
		if firstList[index1][0] <= secondList[index2][0] {
			if firstList[index1][1] < secondList[index2][0] {
				index1++
			} else {
				if firstList[index1][1] <= secondList[index2][1] {
					result = append(result, []int{secondList[index2][0], firstList[index1][1]})
					index1++
				} else {
					result = append(result, []int{secondList[index2][0], secondList[index2][1]})
					index2++
				}
			}
		} else {
			if firstList[index1][0] > secondList[index2][1] {
				index2++
			} else {
				if firstList[index1][1] <= secondList[index2][1] {
					result = append(result, []int{firstList[index1][0], firstList[index1][1]})
					index1++
				} else {
					result = append(result, []int{firstList[index1][0], secondList[index2][1]})
					index2++
				}
			}
		}
	}
	return result
}
