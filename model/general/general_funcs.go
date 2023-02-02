package general

import (
	"math/bits"
	"math/rand"
)

func Choice(nodes []int, k int) []int {
	res := make([]int, 0, k)

	var val int
	for i := 0; i < k; i++ {
		val = rand.Intn(len(nodes))
		res = append(res, nodes[val-1])
	}
	return res
}

func BitLength(num int) int {
	return bits.Len(uint(num))
}

func Contains[T comparable](elems []T, value T) bool {
	for _, item := range elems {
		if item == value {
			return true
		}
	}
	return false
}

func BinarySearchClosest(arr []int, target int, n int) []int {
	left, right := 0, len(arr)-1
	mid := 0
	for left <= right {
		mid = (left + right) / 2
		curNodeId := arr[mid]
		if curNodeId == target {
			//if curNodeId > target-(n/2) && curNodeId < target+(n/2) {
			break
		} else if curNodeId < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	left = mid - n
	if left < 0 {
		left = 0
	}
	right = mid + n
	if right > len(arr) {
		right = len(arr)
	}
	return arr[left:right]
}
