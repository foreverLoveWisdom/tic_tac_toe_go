package main

import (
	"log"
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

func main() {
	board := InitializeBoard()
	log.Println(board)
}
