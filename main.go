package main

import "fmt"

func main() {

	row, col := 8, 8

	var grid [][]int = make([][]int, row)

	for i := range grid {
		grid[i] = make([]int, col)
	}

	for _, row := range grid {
		fmt.Println(row)
	}

}
