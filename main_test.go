package main

import "testing"

func TestInitialBoard(t *testing.T) {
	board := InitializeBoard()
	for row := range [3]int{} {
		for col := range [3]int{} {
			if board[row][col] != ' ' {
				t.Errorf("Expected cell (%d,%d) to be empty, got %c", row, col, board[row][col])
			}
		}
	}
}
