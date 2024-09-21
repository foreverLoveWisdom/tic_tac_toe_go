package main

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	ColorReset     = "\033[0m"
	ColorRed       = "\033[31m"
	ColorBlue      = "\033[34m"
	ColorBold      = "\033[1m"
	ColorDim       = "\033[2m"
	ColorUnderline = "\033[4m"
)

var lastRow, lastCol int = -1, -1

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			// handle the error
		}
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			// handle the error
		}
	default:
		// Unsupported OS; do nothing
	}
}

func ColorizeRune(r rune, highlight bool) string {
	colored := ""
	switch r {
	case 'X':
		colored = ColorRed + string(r) + ColorReset
	case 'O':
		colored = ColorBlue + string(r) + ColorReset
	default:
		colored = " "
	}

	if highlight && r != ' ' {
		return ColorBold + colored + ColorReset
	}
	return colored
}

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

	// Column headers
	sb.WriteString("   1   2   3\n")

	for row := range [3]struct{}{} {
		// Row number
		sb.WriteString(strconv.Itoa(row+1) + "  ")
		for col := range [3]struct{}{} {
			highlight := (row == lastRow && col == lastCol)
			sb.WriteString(ColorizeRune(board[row][col], highlight))
			const maxCol = 2

			if col < maxCol {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString("\n")
		const maxRow = 2

		if row < maxRow {
			sb.WriteString("  ---+---+---\n")
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

func CheckDraw(board [3][3]rune) bool {
	for row := range [3]struct{}{} {
		for col := range [3]struct{}{} {
			if board[row][col] == ' ' {
				return false
			}
		}
	}

	return true
}

func main() {
	board := InitializeBoard()
	currentPlayer := 'X'
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		os.Stdout.WriteString(DisplayBoard(board) + "\n")
		os.Stdout.WriteString("Player " + string(currentPlayer) + ", enter your move (row and column: 1 1): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			os.Stdout.WriteString("Error reading input, please try again.\n")
			continue
		}

		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")
		const expectedParts = 2

		if len(parts) != expectedParts {
			os.Stdout.WriteString("Invalid input format. Please enter two numbers separated by a space.\n")
			continue
		}

		row, err1 := strconv.Atoi(parts[0])
		col, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			os.Stdout.WriteString("Invalid numbers. Please enter integers between 1 and 3.\n")
			continue
		}

		row--
		col--

		if !IsValidMove(board, row, col) {
			os.Stdout.WriteString("Invalid move. Cell is either occupied or out of range.\n")
			continue
		}

		board, err = ApplyMove(board, row, col, currentPlayer)
		if err != nil {
			os.Stdout.WriteString("Error applying move. Please try again.\n")
			continue
		}

		if CheckWin(board, currentPlayer) {
			clearScreen()
			os.Stdout.WriteString(DisplayBoard(board) + "\n")
			os.Stdout.WriteString("Player " + string(currentPlayer) + " wins!\n")
			break
		}

		if CheckDraw(board) {
			clearScreen()
			os.Stdout.WriteString(DisplayBoard(board) + "\n")
			os.Stdout.WriteString("It's a draw!\n")
			break
		}

		if currentPlayer == 'X' {
			currentPlayer = 'O'
		} else {
			currentPlayer = 'X'
		}
	}
}
