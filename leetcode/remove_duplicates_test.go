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

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/explore/interview/card/top-interview-questions-easy/1/array/21/
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
package leetcode_test

import (
	"reflect"
	"testing"
)

/* solution */
func removeDuplicates(nums []int) []int {
	//fmt.Printf("\n nums：%d",nums)
	uniqueCount := 0
	repeatCount := 0
	for i,_ := range nums{
		//fmt.Printf("\n =====i：%d",i)
		if i == 0 {
			uniqueCount++
			continue
		}
		unique := true
		// 向前遍历 是否已存在
		for j := i - 1 ; j >= 0 ; j-- {
			//fmt.Printf(", j：%d",j)
			if nums[j] == nums[i] {
				unique = false
				repeatCount++
				break
			}
		}
		if unique {
			//fmt.Printf("\n 新")
			uniqueCount++
			nums[i], nums[uniqueCount-1] = nums[uniqueCount-1], nums[i]
		}
		//fmt.Printf("\n nums：%d",nums)
	}
	//fmt.Printf("\n nums：%d",nums)
	//fmt.Printf("\n\t uniqueCount：%d, repeatCount：%d",uniqueCount, repeatCount)
	return nums[:uniqueCount]
}

/* test */
func TestRemoveDuplicates(t *testing.T) {
	type problem struct {
		nums []int
		expect []int
	}
	problemList := []problem{
		{
			[]int{},
			[]int{},
		},
		{
			[]int{1, 1, 2},
			[]int{1, 2},
		},
		{
			[]int{0,0,1,1,1,2,2,3,3,4},
			[]int{0,1,2,3,4},
		},
	}
	for _, item := range problemList {
		expect := removeDuplicates(item.nums)
		if !reflect.DeepEqual(expect, item.expect) {
			t.Errorf("\n expect：%d\n actual：%d", item.expect, expect)
		}
	}
}