/*
一个机器人位于一个 m x n 网格的左上角 （起始点在下图中标记为“Start” ）。
机器人每次只能向下或者向右移动一步。机器人试图达到网格的右下角（在下图中标记为“Finish”）。
问总共有多少条不同的路径？

示例 1:
输入: m = 3, n = 2
输出: 3
解释:
从左上角开始，总共有 3 条路径可以到达右下角。
1. 向右 -> 向右 -> 向下
2. 向右 -> 向下 -> 向右
3. 向下 -> 向右 -> 向右

示例 2:
输入: m = 7, n = 3
输出: 28

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/unique-paths
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
package leetcode_test

import (
	"math/big"
	"testing"
)

/* solution */
// method 1 : 枚举路径法 (数据量过大,舍弃)
func uniquePathsByEnum(m int, n int) (wayCount int) {
	return
}
// method 2 : 排列组合法
func uniquePathsByCombination(m int, n int) int {
	if m == 1 || n == 1 {
		return 1
	}
	stepToRight := m - 1
	stepToBottom := n - 1
	stepAll := stepToRight + stepToBottom

	allOrder := factorial(stepAll)
	rightOrder := factorial(stepToRight)
	bottomOrder := factorial(stepToBottom)
	numeratorOfwayCount := allOrder
	var denominatorOfwayCount big.Int
	denominatorOfwayCount.Mul(&bottomOrder, &rightOrder)
	var wayCount big.Int
	wayCount.Div(&numeratorOfwayCount, &denominatorOfwayCount)
	return int(wayCount.Uint64())
}
func factorial (num int) (value big.Int) {
	// 阶乘 计算溢出问题 通过big库解锁
	value.SetUint64(1)
	for i := 1 ; i <= num ; i++ {
		n := int64(i)
		value.Mul(big.NewInt(n), &value)
	}
	return value
}
/* test */
func TestUniquePathsByCombination(t *testing.T) {
	type matrix struct {
		m int
		n int
		expect int
	}
	matrixList := []matrix{
		{
			3,
			2,
			3,
		},
		{
			7,
			3,
			28,
		},
		{
			23,
			12,
			193536720,
		},
	}
	for _, item := range matrixList {
		expect := uniquePathsByCombination(item.m, item.n)
		if expect != item.expect {
			t.Errorf("\nexpect：%d\nactual：%d", item.expect, expect)
		}
	}
}