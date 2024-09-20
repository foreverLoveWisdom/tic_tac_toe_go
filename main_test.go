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

func TestDisplayBoard(t *testing.T) {
	board := InitializeBoard()
	board[0][0] = 'X'
	board[1][1] = 'O'
	expectedOutput := " X |   |  \n---+---+---\n   | O |  \n---+---+---\n   |   |  \n"

	got := DisplayBoard(board)
	if got != expectedOutput {
		t.Errorf("DisplayBoard() =\n%s\nExpected:\n%s", got, expectedOutput)
	}
}
