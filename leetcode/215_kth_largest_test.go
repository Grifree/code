/*
215. 数组中的第K个最大元素
在未排序的数组中找到第 k 个最大的元素。请注意，你需要找的是数组排序后的第 k 个最大的元素，而不是第 k 个不同的元素。

示例 1:
输入: [3,2,1,5,6,4] 和 k = 2
输出: 5

示例 2:
输入: [3,2,3,1,2,4,5,5,6] 和 k = 4
输出: 4

说明:
你可以假设 k 总是有效的，且 1 ≤ k ≤ 数组的长度。
*/
package leetcode_test

import (
	"testing"
)

// solution
func findKthLargest(nums []int, k int) int {
	listDesc := sortByDesc(nums)
	kth := listDesc[k - 1]
	return kth
}
func sortByDesc (list []int) (listDesc []int) {
	listDesc = list
	for i := 0 ; i < len(listDesc) ; i++ {
		for j := i + 1 ; j < len(listDesc) ; j++ {
			if listDesc[j] > listDesc[i] {
				listDesc[j],listDesc[i] = listDesc[i],listDesc[j]
			}
		}
	}
	return listDesc
}

// test
func TestFindKthLargest(t *testing.T) {
	type problem struct {
		nums []int
		k int
		expect int
	}
	problemList := []problem{
		{
			[]int{3,2,1,5,6,4},
			2,
			5,
		},
		{
			[]int{3,2,3,1,2,4,5,5,6},
			4,
			4,
		},
	}
	for _, item := range problemList {
		expect := findKthLargest(item.nums, item.k)
		if expect != item.expect {
			t.Errorf("expect：%d\nactual：%d", item.expect, expect)
		}
	}
}