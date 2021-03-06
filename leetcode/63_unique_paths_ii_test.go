/* 63. 不同路径 II
一个机器人位于一个 m x n 网格的左上角 （起始点在下图中标记为“Start” ）。
机器人每次只能向下或者向右移动一步。机器人试图达到网格的右下角（在下图中标记为“Finish”）。
现在考虑网格中有障碍物。那么从左上角到右下角将会有多少条不同的路径？

网格中的障碍物和空位置分别用 1 和 0 来表示。

说明：m 和 n 的值均不超过 100。

示例 1:
输入:
[
  [0,0,0],
  [0,1,0],
  [0,0,0]
]
输出: 2
解释:
3x3 网格的正中间有一个障碍物。
从左上角到右下角一共有 2 条不同的路径：
1. 向右 -> 向右 -> 向下 -> 向下
2. 向下 -> 向下 -> 向右 -> 向右

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/unique-paths-ii
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
package leetcode_test

import (
	"fmt"
	"testing"
)

/* solution */
// method 1 : 递归法 (大矩阵性能极差)
type ObstacleGrid [][]int
func uniquePathsWithObstaclesByRecursion(obstacleGrid ObstacleGrid) int {
	y := len(obstacleGrid) - 1
	x := len(obstacleGrid[0]) - 1
	//fmt.Printf("\n ------------ obstacleGrid：%d，%d", x, y)
	count := dp(obstacleGrid, Vector{x,y})
	return count
}
type Vector struct {
	x int
	y int
}
func dp(obstacleGrid ObstacleGrid, vector Vector) int {
	//fmt.Printf("\n vector：%d，%d", vector.x, vector.y)
	// 判断顺序不能变
	if vector.y < 0 || vector.x < 0 {
		//fmt.Printf("超出")
		return 0
	}
	if obstacleGrid[vector.y][vector.x] == 1 {
		//fmt.Printf("障碍")
		return 0
	}
	if vector.y == 0 && vector.x == 0 {
		//fmt.Printf("起点")
		return 1
	}
	return dp(obstacleGrid, Vector{vector.x - 1,vector.y}) + dp(obstacleGrid, Vector{vector.x, vector.y - 1})
}

// method 2 : 遍历计数法
func uniquePathsWithObstaclesByCount(obstacleGrid ObstacleGrid) int {
	row := len(obstacleGrid) // y
	col := len(obstacleGrid[0]) // x
	for y := 0 ; y < row ; y++ {
		for x := 0 ; x < col ; x++ {
			// 有障碍物
			if obstacleGrid[y][x] == 1 {
				obstacleGrid[y][x] = 0
				continue
			}
			// 无障碍物 计数
			if x == 0 && y == 0 { /*起点*/
				obstacleGrid[y][x] = 1
				continue;
			}
			/*其余点*/
			top := 0
			if x - 1 >= 0 {
				top = obstacleGrid[y][x-1]
			}
			left := 0
			if y - 1 >= 0 {
				left = obstacleGrid[y-1][x]
			}
			obstacleGrid[y][x] = top + left
		}
	}
	return obstacleGrid[row-1][col-1]
}

/* test */
func TestUniquePathsWithObstacles(t *testing.T) {
	type problem struct {
		obstacleGrid ObstacleGrid
		expect int
	}
	problemList := []problem{
		{
			ObstacleGrid{
				{0,0,0},
				{0,1,0},
				{0,0,0},
			},
			2,
		},
		{
			ObstacleGrid{
				{0,0,0},
				{0,0,0},
				{0,0,0},
			},
			6,
		},
		{
			ObstacleGrid{
				{0,0,0},
			},
			1,
		},
		{
			ObstacleGrid{
				{0,0,0},
				{0,0,1},
				{1,0,0},
			},
			2,
		},
		{
			ObstacleGrid{
				{0,1},
				{0,0},
			},
			1,
		},
		{
			ObstacleGrid{
				{0,0,1,0},
				{0,0,1,0},
				{1,0,0,0},
				{0,0,0,0},
			},
			6,
		},
		{
			ObstacleGrid{
				{1},
			},
			0,
		},
		{
			ObstacleGrid{
				{0},
			},
			1,
		},
		{
			ObstacleGrid{
				{0,0,0,0,0,0,1,0,0,0,0,1,0,0,0,0,0,0,1,0,0,0,0,0,0,1,0,1,0},
				{1,0,0,0,0,1,0,0,1,0,0,0,1,0,1,0,0,0,0,0,0,0,1,1,0,0,0,0,1},
				{0,0,1,0,0,1,0,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,1,0,0},
				{0,0,0,0,1,0,1,1,0,0,0,0,0,0,0,0,0,1,0,1,0,1,0,0,0,0,0,0,0},
				{0,0,0,1,0,0,1,1,0,0,1,0,0,0,1,0,0,0,0,0,0,1,0,0,0,0,0,0,0},
				{0,0,1,0,1,0,1,0,1,0,1,0,0,1,0,0,0,1,0,0,1,0,0,1,0,0,1,0,0},
				{1,0,1,0,0,0,0,0,0,1,0,0,0,0,1,0,0,0,0,0,0,0,1,1,0,1,0,1,0},
				{0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0},
				{1,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,0,1,0,0,0,0,1,0,0,1,1,0,1},
				{0,0,0,1,0,0,0,1,1,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,1,0,0,0},
				{0,0,1,1,0,0,0,1,1,1,0,0,0,1,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0},
				{0,1,0,0,0,0,0,0,1,0,0,0,0,1,0,0,0,0,0,0,0,1,0,0,0,0,1,0,0},
				{0,0,1,0,0,0,1,0,0,0,0,0,1,1,0,1,0,0,0,0,0,1,0,0,0,0,0,0,0},
				{0,0,0,1,1,0,0,0,0,1,1,0,0,0,0,0,0,0,0,0,1,0,1,0,0,1,0,0,0},
				{0,0,0,1,0,0,0,0,1,0,1,1,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,1,0},
				{0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
				{0,0,0,0,1,0,1,0,0,0,0,1,0,0,0,1,0,0,0,1,1,1,1,0,0,0,1,0,1},
				{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,1},
				{0,0,0,0,0,1,0,1,0,0,0,1,0,0,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0},
				{0,0,0,0,0,0,0,0,0,0,1,1,0,1,1,0,0,0,0,1,1,0,0,0,0,0,0,1,0},
				{0,0,1,0,0,0,0,0,1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,1},
				{0,1,1,1,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0},
				{0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,1,0,0,0,0,0,1,1,0,1,0,0},
				{0,0,0,0,0,1,1,0,0,0,0,1,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0},
				{0,0,0,0,1,1,0,1,0,1,1,1,1,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,1},
				{0,1,0,0,0,0,0,1,1,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,1,0,1,0},
				{0,1,0,0,0,0,0,1,0,1,1,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0},
				{0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,1,0,0,1,0,0,0,0,0,0,0,1,0},
			},
			2768280,
		},
	}
	for _, item := range problemList {
		expect := uniquePathsWithObstaclesByCount(item.obstacleGrid)
		if expect != item.expect {
			t.Errorf("\n expect：%d\n actual：%d", item.expect, expect)
		}
	}
}