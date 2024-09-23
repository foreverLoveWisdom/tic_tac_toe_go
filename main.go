package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
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

// ---High-level Game Flow---.
func main() {
	// Remove timestamp from log messages.
	log.SetFlags(0)

	for {
		board, currentPlayer, reader := SetupNewGame()
		RunGame(board, currentPlayer, reader)

		if !PromptRestart(reader) {
			log.Println("Thank you for playing! Goodbye.")

			break
		}
	}
}

func SetupNewGame() ([3][3]rune, rune, *bufio.Reader) {
	board := initializeBoard()
	player := 'X'
	reader := bufio.NewReader(os.Stdin)

	printWelcomeMessage()

	return board, player, reader
}

func RunGame(board [3][3]rune, currentPlayer rune, reader *bufio.Reader) ([3][3]rune, rune) {
	var valid bool

	for {
		printBoard(board)
		board, currentPlayer, valid = processPlayerMove(board, currentPlayer, reader)

		if !valid {
			continue
		}

		gameEnded, message := evaluateGameStatus(board, currentPlayer)
		if gameEnded {
			printBoard(board)
			log.Println(message)

			break
		}

		currentPlayer = switchPlayer(currentPlayer)
	}

	return board, currentPlayer
}

func PromptRestart(reader *bufio.Reader) bool {
	log.Println("Game over! Would you like to play again? (y/n): ")

	input, _ := reader.ReadString('\n')

	return strings.TrimSpace(input) == "y"
}

// ---Game Logic---.
func evaluateGameStatus(board [3][3]rune, currentPlayer rune) (bool, string) {
	if checkWin(board, currentPlayer) {
		return true, fmt.Sprintf("Player %c wins!", currentPlayer)
	}

	if checkDraw(board) {
		return true, "It's a draw!"
	}

	return false, ""
}

func processPlayerMove(board [3][3]rune, currentPlayer rune, reader *bufio.Reader) ([3][3]rune, rune, bool) {
	row, col, err := getPlayerMove(reader, currentPlayer)
	if err != nil {
		log.Println(err)
		return board, currentPlayer, false
	}

	board, err = executeMove(board, row, col, currentPlayer)
	if err != nil {
		log.Println("Invalid move. Try again.")
		return board, currentPlayer, false
	}

	return board, currentPlayer, true
}

func switchPlayer(currentPlayer rune) rune {
	if currentPlayer == 'X' {
		return 'O'
	}

	return 'X'
}

// --- Board Management (Board Display & Move Execution) ---.
func printBoard(board [3][3]rune) {
	os.Stdout.WriteString(renderBoard(board) + "\n")
}

func renderBoard(board [3][3]rune) string {
	const maxCol = 2

	const maxRow = 2

	var sb strings.Builder

	sb.WriteString("   1   2   3\n")

	for row := range [3]struct{}{} {
		sb.WriteString(strconv.Itoa(row+1) + "  ")

		for col := range [3]struct{}{} {
			sb.WriteString(colorizeRune(board[row][col]))

			if col < maxCol {
				sb.WriteString(" | ")
			}
		}

		sb.WriteString("\n")

		if row < maxRow {
			sb.WriteString("  ---+---+---\n")
		}
	}

	return sb.String()
}

func executeMove(board [3][3]rune, row, col int, currentPlayer rune) ([3][3]rune, error) {
	if !isValidMove(board, row, col) {
		return board, errors.New("invalid move")
	}

	return applyMove(board, row, col, currentPlayer)
}

func applyMove(board [3][3]rune, currentRow int, currentCol int, player rune) ([3][3]rune, error) {
	if player != 'X' && player != 'O' {
		return board, errors.New("invalid player")
	}

	if !isValidMove(board, currentRow, currentCol) {
		return board, errors.New("invalid move")
	}

	newBoard := board
	newBoard[currentRow][currentCol] = player

	return newBoard, nil
}

func isValidMove(board [3][3]rune, currentRow int, currentCol int) bool {
	if currentRow < 0 || currentRow >= 3 || currentCol < 0 || currentCol >= 3 {
		return false
	}

	return board[currentRow][currentCol] == ' '
}

func initializeBoard() [3][3]rune {
	var board [3][3]rune

	for row := range [3]struct{}{} {
		for col := range [3]struct{}{} {
			board[row][col] = ' '
		}
	}

	return board
}

// ---Input handling---.
func getPlayerMove(reader *bufio.Reader, currentPlayer rune) (int, int, error) {
	input := promptPlayerMove(reader, currentPlayer)
	return parseMove(input)
}

func promptPlayerMove(reader *bufio.Reader, currentPlayer rune) string {
	message := "Player %c, it's your turn! Please enter your move as guided above: "
	log.Printf(message, currentPlayer)

	input, _ := reader.ReadString('\n')

	return strings.TrimSpace(input)
}

func parseMove(input string) (int, int, error) {
	const expectedParts = 2

	parts := strings.Split(input, " ")
	if len(parts) != expectedParts {
		return -1, -1, errors.New("invalid input format")
	}

	row, err1 := strconv.Atoi(parts[0])
	col, err2 := strconv.Atoi(parts[1])

	if err1 != nil {
		return -1, -1, errors.New("invalid row number")
	}

	if err2 != nil {
		return -1, -1, errors.New("invalid column number")
	}

	return row - 1, col - 1, nil
}

// ---Win/Draw Condition Checking---.
func checkWin(board [3][3]rune, player rune) bool {
	return checkRows(board, player) || checkColumns(board, player) || checkDiagonals(board, player)
}

func checkDraw(board [3][3]rune) bool {
	for row := range [3]struct{}{} {
		for col := range [3]struct{}{} {
			if board[row][col] == ' ' {
				return false
			}
		}
	}

	return true
}

func checkRows(board [3][3]rune, player rune) bool {
	for row := range [3]struct{}{} {
		if isRowWin(board, player, row) {
			return true
		}
	}

	return false
}

func checkColumns(board [3][3]rune, player rune) bool {
	for col := range [3]struct{}{} {
		if isColumnWin(board, player, col) {
			return true
		}
	}

	return false
}

func checkDiagonals(board [3][3]rune, player rune) bool {
	return isMainDiagonalWin(board, player) || isAntiDiagonalWin(board, player)
}

func isRowWin(board [3][3]rune, player rune, row int) bool {
	return board[row][0] == player && board[row][1] == player && board[row][2] == player
}

func isColumnWin(board [3][3]rune, player rune, col int) bool {
	return board[0][col] == player && board[1][col] == player && board[2][col] == player
}

func isMainDiagonalWin(board [3][3]rune, player rune) bool {
	return board[0][0] == player && board[1][1] == player && board[2][2] == player
}

func isAntiDiagonalWin(board [3][3]rune, player rune) bool {
	return board[0][2] == player && board[1][1] == player && board[2][0] == player
}

// ---Utility Functions---.
func colorizeRune(r rune) string {
	switch r {
	case 'X':
		return ColorBold + ColorRed + string(r) + ColorReset
	case 'O':
		return ColorBold + ColorBlue + string(r) + ColorReset
	default:
		return " "
	}
}

func printWelcomeMessage() {
	log.Println("Welcome to Tic-Tac-Toe!")
	log.Println("Players take turns to place their mark (X or O) on the board.")
	log.Println("Enter your move as 'row column' (e.g., '1 1' for top-left corner).")
	log.Println("Note: There is a white space between the row and column.")
	log.Println("Rows and columns are numbered from 1 to 3.")
}
