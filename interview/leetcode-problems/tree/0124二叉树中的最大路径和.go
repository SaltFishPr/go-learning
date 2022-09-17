// @file: 0124二叉树中的最大路径和.go
// @description: 给定一个非空二叉树，返回其最大路径和。
// @author: SaltFish
// @date: 2020/09/05

package tree

/*
本题中，路径被定义为一条从树中任意节点出发，沿父节点-子节点连接，达到任意节点的序列。该路径至少包含一个节点，且不一定经过根节点。

示例 1：
输入：[1,2,3]

       1
      / \
     2   3

输出：6

示例 2：
输入：[-10,9,20,null,null,15,7]

   -10
   / \
  9  20
    /  \
   15   7

输出：42
*/

import "math"

func maxPathSum(root *TreeNode) int {
	maxSum := math.MinInt32
	var maxGain func(*TreeNode) int // 先声明函数
	maxGain = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		leftGain := max(maxGain(node.Left), 0)
		rightGain := max(maxGain(node.Right), 0)
		priceNewPath := node.Val + leftGain + rightGain
		maxSum = max(maxSum, priceNewPath)
		return node.Val + max(leftGain, rightGain)
	}
	maxGain(root)
	return maxSum
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
