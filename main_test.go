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
	board, _, _ := SetupNewGame()
	board, _, valid := processPlayerMove(board, 'X', createReader("1 1\n"))

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

// TestSwitchPlayer checks if players are switched correctly.
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

// Helper function to simulate user input.
func createReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}
