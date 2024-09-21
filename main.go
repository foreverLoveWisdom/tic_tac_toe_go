package main

import (
	"bufio"
	"errors"
	"log"
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

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to clear screen: %v\n", err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to clear screen: %v\n", err)
		}
	default:
		log.Printf("Unsupported operating system: %s\n", runtime.GOOS)
	}
}

func ColorizeRune(r rune) string {
	switch r {
	case 'X':
		return ColorBold + ColorRed + string(r) + ColorReset
	case 'O':
		return ColorBold + ColorBlue + string(r) + ColorReset
	default:
		return " "
	}
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

func DisplayBoard(board [3][3]rune, lastRow int, lastCol int) string {
	var sb strings.Builder

	sb.WriteString("   1   2   3\n")

	for row := range [3]struct{}{} {
		sb.WriteString(strconv.Itoa(row+1) + "  ")
		for col := range [3]struct{}{} {
			sb.WriteString(ColorizeRune(board[row][col]))
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

func printWelcomeMessage() {
	log.Println("Welcome to Tic-Tac-Toe!")
	log.Println("Players take turns to place their mark (X or O) on the board.")
	log.Println("Enter your move as 'row column' (e.g., '1 1' for top-left corner).")
	log.Println("Rows and columns are numbered from 1 to 3.")
}

func promptPlayerMove(reader *bufio.Reader, currentPlayer rune) string {
	message := "Player %c, it's your turn! Please enter your move as guided above: "
	log.Printf(message, currentPlayer)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func parseMove(input string) (int, int, error) {
	parts := strings.Split(input, " ")
	const expectedParts = 2

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

func switchPlayer(currentPlayer rune) rune {
	if currentPlayer == 'X' {
		return 'O'
	}

	return 'X'
}

func main() {
	for {
		board := InitializeBoard()
		lastRow, lastCol := -1, -1
		currentPlayer := 'X'
		reader := bufio.NewReader(os.Stdin)

		printWelcomeMessage()

		for {
			os.Stdout.WriteString(DisplayBoard(board, lastRow, lastCol) + "\n")
			input := promptPlayerMove(reader, currentPlayer)

			row, col, err := parseMove(input)
			if err != nil {
				clearScreen()
				log.Println(err)
				continue
			}

			if !IsValidMove(board, row, col) {
				clearScreen()
				log.Println("Invalid move. Cell is either occupied or out of range.")
				continue
			}

			board, _ = ApplyMove(board, row, col, currentPlayer)
			lastRow, lastCol = row, col

			if CheckWin(board, currentPlayer) {
				clearScreen()
				os.Stdout.WriteString(DisplayBoard(board, lastRow, lastCol) + "\n")
				log.Printf("Player %c wins!\n", currentPlayer)
				break
			}

			if CheckDraw(board) {
				clearScreen()
				os.Stdout.WriteString(DisplayBoard(board, lastRow, lastCol) + "\n")
				log.Println("It's a draw!")
				break
			}

			currentPlayer = switchPlayer(currentPlayer)
		}

		log.Println("Game over! Would you like to play again? (y/n): ")
		restartInput, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(restartInput)) != "y" {
			clearScreen()
			log.Println("Thank you for playing! Goodbye.")
			break
		}
	}
}
