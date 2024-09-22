package main

import (
	"regexp"
	"testing"
)

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

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(str string) string {
	return ansi.ReplaceAllString(str, "")
}

func TestDisplayBoard(t *testing.T) {
	board := InitializeBoard()
	board[0][0] = 'X'
	board[1][1] = 'O'

	expectedOutput := "   1   2   3\n" +
		"1  " + ColorRed + "X" + ColorReset + " |   |  \n" +
		"  ---+---+---\n" +
		"2    | " + ColorBold + ColorBlue + "O" + ColorReset + " |  \n" +
		"  ---+---+---\n" +
		"3    |   |  \n"

	got := DisplayBoard(board)
	cleanGot := stripANSI(got)
	cleanExpected := stripANSI(expectedOutput)

	if cleanGot != cleanExpected {
		t.Errorf("DisplayBoard() =\n%q\nExpected:\n%q", got, expectedOutput)
	}
}

func isTargetCell(i, j, row, col int) bool {
	return i == row && j == col
}

func isCellChanged(oldBoard, newBoard [3][3]rune, i, j int) bool {
	return newBoard[i][j] != oldBoard[i][j]
}

func verifyUnchangedCells(t *testing.T, oldBoard [3][3]rune, newBoard [3][3]rune, row int, col int) {
	for i := range [3]int{} {
		for j := range [3]int{} {
			if !isTargetCell(i, j, row, col) {
				if isCellChanged(oldBoard, newBoard, i, j) {
					t.Errorf("ApplyMove() modified unexpected cell at (%d, %d)", i, j)
				}
			}
		}
	}
}

func TestIsValidMove(t *testing.T) {
	board := InitializeBoard()
	board[1][1] = 'O'

	tests := []struct {
		name     string
		row, col int
		expected bool
	}{
		{"Out-of-bounds row (negative)", -1, 0, false},
		{"Out-of-bounds column (negative)", 0, -1, false},
		{"Out-of-bounds row (too large)", 3, 0, false},
		{"Out-of-bounds column (too large)", 0, 3, false},
		{"Out-of-bounds both", 3, 3, false},
		{"Edge case: last row, first column", 2, 0, true},
		{"Edge case: first row, last column", 0, 2, true},
		{"Occupied cell", 1, 1, false},
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

func TestApplyMove(t *testing.T) {
	board := InitializeBoard()

	tests := []struct {
		name     string
		row, col int
		player   rune
		expected rune
		wantErr  bool
	}{
		{"Valid move", 1, 1, 'X', 'X', false},
		{"Out-of-bounds row (negative)", -1, 0, 'O', ' ', true},
		{"Out-of-bounds row (too large)", 3, 0, 'O', ' ', true},
		{"Out-of-bounds column (negative)", 0, -1, 'O', ' ', true},
		{"Out-of-bounds column (too large)", 0, 3, 'O', ' ', true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldBoard := board
			newBoard, err := ApplyMove(board, tt.row, tt.col, tt.player)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ApplyMove() expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("ApplyMove() unexpected error: %v", err)
				return
			}

			if newBoard[tt.row][tt.col] != tt.expected {
				t.Errorf("ApplyMove() = %c at (%d, %d); want %c", newBoard[tt.row][tt.col], tt.row, tt.col, tt.expected)
			}

			verifyUnchangedCells(t, oldBoard, newBoard, tt.row, tt.col)
		})
	}
}

func TestCheckWin(t *testing.T) {
	tests := []struct {
		name     string
		board    [3][3]rune
		player   rune
		expected bool
	}{
		{
			name: "Complete row win",
			board: [3][3]rune{
				{'X', 'X', 'X'},
				{' ', ' ', ' '},
				{' ', ' ', ' '},
			},
			player:   'X',
			expected: true,
		},
		{
			name: "Complete column win",
			board: [3][3]rune{
				{'O', ' ', ' '},
				{'O', ' ', ' '},
				{'O', ' ', ' '},
			},
			player:   'O',
			expected: true,
		},
		{
			name: "Complete diagonal win",
			board: [3][3]rune{
				{' ', ' ', 'X'},
				{' ', 'X', ' '},
				{'X', ' ', ' '},
			},
			player:   'X',
			expected: true,
		},
		{
			name: "Incomplete row win",
			board: [3][3]rune{
				{'X', 'X', ' '},
				{' ', ' ', ' '},
				{' ', ' ', ' '},
			},
			player:   'X',
			expected: false,
		},
		{
			name: "Incomplete column win",
			board: [3][3]rune{
				{'O', ' ', ' '},
				{'O', ' ', ' '},
				{' ', ' ', ' '},
			},
			player:   'O',
			expected: false,
		},
		{
			name: "Incomplete diagonal win",
			board: [3][3]rune{
				{' ', ' ', 'X'},
				{' ', 'X', ' '},
				{' ', ' ', ' '},
			},
			player:   'X',
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckWin(tt.board, tt.player)
			if got != tt.expected {
				t.Errorf("CheckWin() = %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckDraw(t *testing.T) {
	var board [3][3]rune
	board = [3][3]rune{
		{'X', 'O', 'X'},
		{'X', 'O', 'O'},
		{'O', 'X', 'X'},
	}

	if !CheckDraw(board) {
		t.Errorf("Expected the game to be a draw")
	}

	board = InitializeBoard()
	if CheckDraw(board) {
		t.Errorf("Expected the game not to be a draw when the board is empty")
	}

	board = [3][3]rune{
		{'X', 'O', 'X'},
		{'X', 'O', 'O'},
		{'O', 'X', ' '},
	}
	if CheckDraw(board) {
		t.Errorf("Expected the game not to be a draw when there are empty cells")
	}
}
