package main

import "log"

func InitializeBoard() [3][3]rune {
	var board [3][3]rune

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			board[row][col] = ' '
		}

	}
	return board
}

func main() {
	board := InitializeBoard()
	log.Println(board)
}
