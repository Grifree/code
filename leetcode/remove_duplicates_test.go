/*
给定一个排序数组，你需要在 原地 删除重复出现的元素，使得每个元素只出现一次，返回移除后数组的新长度。
不要使用额外的数组空间，你必须在 原地 修改输入数组 并在使用 O(1) 额外空间的条件下完成。

示例 1:
给定数组 nums = [1,1,2],
函数应该返回新的长度 2, 并且原数组 nums 的前两个元素被修改为 1, 2。
你不需要考虑数组中超出新长度后面的元素。

示例 2:
给定 nums = [0,0,1,1,1,2,2,3,3,4],
函数应该返回新的长度 5, 并且原数组 nums 的前五个元素被修改为 0, 1, 2, 3, 4。
你不需要考虑数组中超出新长度后面的元素。
*/
package leetcode_test

import (
	"reflect"
	"testing"
)

func removeDuplicates(nums []int) []int {
	uniqueIndex := 0
	repeatCount := 0
	for i,_ := range nums {
		//fmt.Printf("\n =====i：%d",i)
		if i == 0 {
			continue
		}
		exist := false
		// 向前遍历 是否已存在
		for j := i - 1 ; j >= 0 ; j-- {
			//fmt.Printf("\n j：%d",j)
			if nums[j] == nums[i] {
				exist = true
				uniqueIndex = i
				//fmt.Printf("\n exist：%d",uniqueIndex)
				break
			}
		}
		if exist && (repeatCount + uniqueIndex < len(nums)) {
			repeatCount += 1
			//fmt.Printf("\n repeatCount：%d",repeatCount)
			//fmt.Printf("\n len(nums)：%d",len(nums))
			//fmt.Printf("\n nums[i]：%d",nums[i])
			//fmt.Printf("\n nums[len(nums) - repeatCount]：%d",nums[len(nums) - repeatCount])
			nums[i], nums[len(nums) - repeatCount] = nums[len(nums) - repeatCount], nums[i]
		}
		//fmt.Printf("\n nums：%d",nums)
	}
	//fmt.Printf("\n nums：%d",nums)
	//fmt.Printf("\n uniqueIndex：%d",uniqueIndex)
	//fmt.Printf("\n nums[:uniqueIndex]：%d",nums[:uniqueIndex])
	return nums[:uniqueIndex]
}

func TestRemoveDuplicates(t *testing.T) {
	type problem struct {
		nums []int
		expect []int
	}
	problemList := []problem{
		{
			[]int{1, 1, 2},
			[]int{1, 2},
		},
	}
	for _, item := range problemList {
		expect := removeDuplicates(item.nums)
		if !reflect.DeepEqual(expect, item.expect) {
			t.Errorf("\n expect：%d\n actual：%d", item.expect, expect)
		}
	}
}