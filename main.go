package main

import (
	"errors"
	"strings"
)

func InitializeBoard() [3][3]rune {
	var board [3][3]rune

	for row := range [3]struct{}{} {
		for col := range [3]struct{}{} {
			board[row][col] = ' '
		}
	}

	return board
}

func DisplayBoard(board [3][3]rune) string {
	var sb strings.Builder

	for row := range [3]struct{}{} {
		for col := range [3]struct{}{} {
			sb.WriteRune(' ')
			sb.WriteRune(board[row][col])
			const maxCol = 2

			if col < maxCol {
				sb.WriteString(" |")
			}
		}

		sb.WriteString("\n")
		const maxRow = 2

		if row < maxRow {
			sb.WriteString("---+---+---\n")
		}
	}

	return sb.String()
}

func IsValidMove(board [3][3]rune, currentRow int, currentCol int) bool {
	if currentRow < 0 || currentRow >= 3 || currentCol < 0 || currentCol >= 3 {
		return false
	}

	return board[currentRow][currentCol] == ' '
}

func ApplyMove(board [3][3]rune, currentRow int, currentCol int, player rune) ([3][3]rune, error) {
	if player != 'X' && player != 'O' {
		return board, errors.New("invalid player")
	}

	if !IsValidMove(board, currentRow, currentCol) {
		return board, errors.New("invalid move")
	}

	newBoard := board
	newBoard[currentRow][currentCol] = player

	return newBoard, nil
}

func CheckWin(board [3][3]rune, player rune) bool {
	for row := range [3]struct{}{} {
		if board[row][0] == player && board[row][1] == player && board[row][2] == player {
			return true
		}
	}

	for col := range [3]struct{}{} {
		if board[0][col] == player && board[1][col] == player && board[2][col] == player {
			return true
		}
	}

	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true
	}

	return false
}

func main() {
	board := InitializeBoard()
	DisplayBoard(board)
}
