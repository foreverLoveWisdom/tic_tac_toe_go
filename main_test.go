package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestSetupNewGame(t *testing.T) {
	board, player, _ := SetupNewGame()

	for row := range board {
		for col := range board[row] {
			if board[row][col] != ' ' {
				t.Errorf("Expected cell (%d,%d) to be empty, but got %c", row, col, board[row][col])
			}
		}
	}

	if player != 'X' {
		t.Errorf("Expected first player to be 'X', but got %c", player)
	}
}

func TestProcessPlayerMove(t *testing.T) {
	var valid bool

	board, _, _ := SetupNewGame()
	board, _, valid = processPlayerMove(board, 'X', createReader("1 1\n"))

	if !valid {
		t.Error("Expected move to be valid")
	}

	if board[0][0] != 'X' {
		t.Errorf("Expected 'X' at (0,0), but got %c", board[0][0])
	}

	_, _, valid = processPlayerMove(board, 'O', createReader("1 1\n"))
	if valid {
		t.Error("Expected move to be invalid")
	}
}

func TestEvaluateGameStatus(t *testing.T) {
	board := [3][3]rune{
		{'X', 'X', 'X'},
		{' ', 'O', ' '},
		{'O', ' ', ' '},
	}
	gameEnded, message := evaluateGameStatus(board, 'X')

	if !gameEnded || message != "Player X wins!" {
		t.Errorf("Expected 'Player X wins!', got %v, %s", gameEnded, message)
	}

	board = [3][3]rune{
		{'X', 'O', 'X'},
		{'X', 'O', 'O'},
		{'O', 'X', 'X'},
	}
	gameEnded, message = evaluateGameStatus(board, 'X')

	if !gameEnded || message != "It's a draw!" {
		t.Errorf("Expected draw message, but got %v, %s", gameEnded, message)
	}
}

func TestSwitchPlayer(t *testing.T) {
	player := 'X'
	nextPlayer := switchPlayer(player)

	if nextPlayer != 'O' {
		t.Errorf("Expected next player to be 'O', but got %c", nextPlayer)
	}

	player = 'O'
	nextPlayer = switchPlayer(player)

	if nextPlayer != 'X' {
		t.Errorf("Expected next player to be 'X', but got %c", nextPlayer)
	}
}

func TestQuitGame(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Restart with lowercase y", "y\n", true},
		{"Restart with uppercase Y", "Y\n", false},
		{"Do not restart with n", "n\n", false},
		{"Do not restart with invalid input", "invalid\n", false},
		{"Do not restart with empty input", "\n", false},
		{"Restart with extra spaces and y", " y \n", true},
		{"Do not restart with extra spaces and n", " n \n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			bufReader := bufio.NewReader(reader)

			result := PlayAgain(bufReader)

			if result != tt.expected {
				t.Errorf("PlayAgain() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Helper function to simulate user input.
func createReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}
