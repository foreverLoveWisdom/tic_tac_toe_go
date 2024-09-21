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

func TestIsValidMove(t *testing.T) {
	board := InitializeBoard()
	board[0][0] = 'X'
	board[1][1] = 'O'

	tests := []struct {
		name     string
		row, col int
		expected bool
	}{
		{"Occupied cell", 0, 0, false},
		{"Valid unoccupied cell", 0, 1, true},
		{"Out-of-bounds row (negative)", -1, 0, false},
		{"Out-of-bounds column (negative)", 0, -1, false},
		{"Out-of-bounds row (too large)", 3, 0, false},
		{"Out-of-bounds column (too large)", 0, 3, false},
		{"Out-of-bounds both", 3, 3, false},
		{"Edge case: last row, first column", 2, 0, true},
		{"Edge case: first row, last column", 0, 2, true},
		{"Empty cell (occupied by 'O')", 1, 1, false},
		{"Valid corner", 2, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidMove(board, tt.row, tt.col)
			if got != tt.expected {
				t.Errorf("IsValidMove(board, %d, %d) = %v; want %v", tt.row, tt.col, got, tt.expected)
			}
		})
	}

	fullBoard := [3][3]rune{
		{'X', 'O', 'X'},
		{'O', 'X', 'O'},
		{'O', 'X', 'O'},
	}
	if IsValidMove(fullBoard, 1, 1) {
		t.Errorf("IsValidMove on full board returned true, expected false")
	}
}
