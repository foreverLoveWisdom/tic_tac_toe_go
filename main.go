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

func DisplayBoard(board [3][3]rune) string {
	const maxCol = 2

	const maxRow = 2

	var sb strings.Builder

	sb.WriteString("   1   2   3\n")

	for row := range [3]struct{}{} {
		sb.WriteString(strconv.Itoa(row+1) + "  ")

		for col := range [3]struct{}{} {
			sb.WriteString(ColorizeRune(board[row][col]))

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
	return checkRows(board, player) || checkColumns(board, player) || checkDiagonals(board, player)
}

func checkRows(board [3][3]rune, player rune) bool {
	for row := range [3]struct{}{} {
		if isRowWin(board, player, row) {
			return true
		}
	}

	return false
}

func isRowWin(board [3][3]rune, player rune, row int) bool {
	return board[row][0] == player && board[row][1] == player && board[row][2] == player
}

func checkColumns(board [3][3]rune, player rune) bool {
	for col := range [3]struct{}{} {
		if isColumnWin(board, player, col) {
			return true
		}
	}

	return false
}

func isColumnWin(board [3][3]rune, player rune, col int) bool {
	return board[0][col] == player && board[1][col] == player && board[2][col] == player
}

func checkDiagonals(board [3][3]rune, player rune) bool {
	return isMainDiagonalWin(board, player) || isAntiDiagonalWin(board, player)
}

func isMainDiagonalWin(board [3][3]rune, player rune) bool {
	return board[0][0] == player && board[1][1] == player && board[2][2] == player
}

func isAntiDiagonalWin(board [3][3]rune, player rune) bool {
	return board[0][2] == player && board[1][1] == player && board[2][0] == player
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

func initializeGame() ([3][3]rune, rune, *bufio.Reader) {
	board := InitializeBoard()
	player := 'X'
	reader := bufio.NewReader(os.Stdin)

	printWelcomeMessage()

	return board, player, reader
}

func runGameLoop(board [3][3]rune, currentPlayer rune, reader *bufio.Reader) ([3][3]rune, rune) {
	var valid bool

	for {
		refreshBoard(board)
		board, currentPlayer, valid = playRound(board, currentPlayer, reader)

		if !valid {
			clearScreen()

			continue
		}

		if CheckWin(board, switchPlayer(currentPlayer)) {
			refreshBoard(board)
			log.Printf("Player %c wins!\n", switchPlayer(currentPlayer))

			break
		}

		if CheckDraw(board) {
			refreshBoard(board)
			log.Println("It's a draw!")

			break
		}
	}

	return board, currentPlayer
}

func printWelcomeMessage() {
	log.Println("Welcome to Tic-Tac-Toe!")
	log.Println("Players take turns to place their mark (X or O) on the board.")
	log.Println("Enter your move as 'row column' (e.g., '1 1' for top-left corner).")
	log.Println("Rows and columns are numbered from 1 to 3.")
}

func refreshBoard(board [3][3]rune) {
	clearScreen()
	os.Stdout.WriteString(DisplayBoard(board) + "\n")
}

func playRound(board [3][3]rune, currentPlayer rune, reader *bufio.Reader) ([3][3]rune, rune, bool) {
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

	return board, switchPlayer(currentPlayer), true
}

func getPlayerMove(reader *bufio.Reader, currentPlayer rune) (int, int, error) {
	input := promptPlayerMove(reader, currentPlayer)
	return parseMove(input)
}

func executeMove(board [3][3]rune, row, col int, currentPlayer rune) ([3][3]rune, error) {
	if !IsValidMove(board, row, col) {
		return board, errors.New("invalid move")
	}

	return ApplyMove(board, row, col, currentPlayer)
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

func switchPlayer(currentPlayer rune) rune {
	if currentPlayer == 'X' {
		return 'O'
	}

	return 'X'
}

func main() {
	for {
		board, currentPlayer, reader := initializeGame()
		runGameLoop(board, currentPlayer, reader)

		log.Println("Game over! Would you like to play again? (y/n): ")

		restartInput, _ := reader.ReadString('\n')

		if strings.TrimSpace(strings.ToLower(restartInput)) != "y" {
			clearScreen()
			log.Println("Thank you for playing! Goodbye.")

			break
		}
	}
}

// Based on the abstraction layers we've identified, here are some **further refactorings** you can apply to your game. These refactorings will help make your code more modular, maintainable, and scalable. Let's go through them layer by layer:

// ### 1. **Game Setup Layer Refactoring**
//    - **Current Behavior:** The board is initialized, and player input is set up directly in `main`.
//    - **Refactoring Idea:** Encapsulate all the setup logic into a dedicated function, `initializeGame`, to make the `main` function more focused.

//    **Refactored Code:**
//    ```go
//    func initializeGame() ([3][3]rune, rune, *bufio.Reader) {
//        board := InitializeBoard()
//        currentPlayer := 'X'
//        reader := bufio.NewReader(os.Stdin)
//        printWelcomeMessage()
//        return board, currentPlayer, reader
//    }
//    ```

//    - **Benefit:** This removes setup clutter from `main` and makes it easier to modify or extend game initialization later (e.g., adding custom rules or settings).

// ### 2. **Game Loop Layer Refactoring**
//    - **Current Behavior:** The game loop is embedded in `main`, mixing responsibilities like input handling and game state updates.
//    - **Refactoring Idea:** Move the core game loop into its own function `runGameLoop`, focusing `main` on high-level game control (start, end, restart).

//    **Refactored Code:**
//    ```go
//    func runGameLoop(board [3][3]rune, currentPlayer rune, reader *bufio.Reader) ([3][3]rune, rune) {
//        for {
//            refreshBoard(board)
//            board, currentPlayer, valid := playRound(board, currentPlayer, reader)
//            if !valid {
//                clearScreen()
//                continue
//            }

//            if CheckWin(board, switchPlayer(currentPlayer)) {
//                refreshBoard(board)
//                log.Printf("Player %c wins!\n", switchPlayer(currentPlayer))
//                break
//            }

//            if CheckDraw(board) {
//                refreshBoard(board)
//                log.Println("It's a draw!")
//                break
//            }
//        }
//        return board, currentPlayer
//    }
//    ```

//    - **Benefit:** The `main` function will now only call `runGameLoop`, which encapsulates the game-playing logic, making `main` cleaner and more focused on controlling the game flow.

// ### 3. **Turn/Action Layer Refactoring**
//    - **Current Behavior:** The turn logic is handled in `playRound`, but it directly modifies the board.
//    - **Refactoring Idea:** Separate player move input from move validation and application. This makes it easier to modify one part (e.g., changing the move input method) without affecting the others.

//    **Refactored Code:**
//    ```go
//    func getPlayerMove(reader *bufio.Reader, currentPlayer rune) (int, int, error) {
//        input := promptPlayerMove(reader, currentPlayer)
//        return parseMove(input)
//    }

//    func executeMove(board [3][3]rune, row, col int, currentPlayer rune) ([3][3]rune, error) {
//        if !IsValidMove(board, row, col) {
//            return board, errors.New("invalid move")
//        }
//        return ApplyMove(board, row, col, currentPlayer)
//    }
//    ```

//    Now, `playRound` will become:
//    ```go
//    func playRound(board [3][3]rune, currentPlayer rune, reader *bufio.Reader) ([3][3]rune, rune, bool) {
//        row, col, err := getPlayerMove(reader, currentPlayer)
//        if err != nil {
//            log.Println(err)
//            return board, currentPlayer, false
//        }

//        board, err = executeMove(board, row, col, currentPlayer)
//        if err != nil {
//            log.Println("Invalid move. Try again.")
//            return board, currentPlayer, false
//        }

//        return board, switchPlayer(currentPlayer), true
//    }
//    ```

//    - **Benefit:** This separates input logic from game action logic, making both pieces easier to modify and test.

// ### 4. **Win/Draw Condition Layer Refactoring**
//    - **Current Behavior:** Win and draw checks are scattered inside the game loop.
//    - **Refactoring Idea:** Encapsulate win and draw checks in a new function `checkGameEnd`. This function will determine whether the game should end, improving modularity.

//    **Refactored Code:**
//    ```go
//    func checkGameEnd(board [3][3]rune, currentPlayer rune) (bool, string) {
//        if CheckWin(board, switchPlayer(currentPlayer)) {
//            return true, fmt.Sprintf("Player %c wins!\n", switchPlayer(currentPlayer))
//        }
//        if CheckDraw(board) {
//            return true, "It's a draw!\n"
//        }
//        return false, ""
//    }
//    ```

//    In your game loop, you can now replace the win/draw checks with:
//    ```go
//    ended, message := checkGameEnd(board, currentPlayer)
//    if ended {
//        refreshBoard(board)
//        log.Println(message)
//        break
//    }
//    ```

//    - **Benefit:** This makes the win/draw conditions more centralized and easier to modify (e.g., adding new rules for winning or extending the grid size).

// ### 5. **Game End/Restart Layer Refactoring**
//    - **Current Behavior:** Game restart logic is embedded in `main`.
//    - **Refactoring Idea:** Move game-ending and restarting logic into a separate function `promptRestart`.

//    **Refactored Code:**
//    ```go
//    func promptRestart(reader *bufio.Reader) bool {
//        log.Println("Game over! Would you like to play again? (y/n): ")
//        restartInput, _ := reader.ReadString('\n')
//        return strings.TrimSpace(strings.ToLower(restartInput)) == "y"
//    }
//    ```

//    Now, replace the restart logic in `main` with:
//    ```go
//    if !promptRestart(reader) {
//        clearScreen()
//        log.Println("Thank you for playing! Goodbye.")
//        break
//    }
//    ```

//    - **Benefit:** This separates the game loop from the restart logic, making it easier to add new features (like saving the game state) or change the game flow.

// ### 6. **Input/Output Layer Refactoring**
//    - **Current Behavior:** Input/output (reading moves, displaying boards) is mixed in various functions.
//    - **Refactoring Idea:** Abstract the input/output operations into separate functions, allowing more flexible UI handling (e.g., switching to a graphical interface or network-based game).

//    **Refactored Code:**
//    ```go
//    func readInput(reader *bufio.Reader) string {
//        input, _ := reader.ReadString('\n')
//        return strings.TrimSpace(input)
//    }

//    func displayMessage(message string) {
//        os.Stdout.WriteString(message + "\n")
//    }
//    ```

//    Replace direct I/O calls in your code with these functions:
//    ```go
//    row, col, err := parseMove(readInput(reader))
//    displayMessage("Player X wins!")
//    ```

//    - **Benefit:** This refactoring abstracts I/O operations, making it easier to swap out terminal-based input/output with another form of UI without changing core game logic.

// ### 7. **Encapsulate Game Management**
//    - **Current Behavior:** `main` is managing multiple responsibilities (initialization, game loop, restarts).
//    - **Refactoring Idea:** Move overall game management into a `TicTacToeGame` struct that handles the state, players, and the game loop.

//    **Refactored Code:**
//    ```go
//    type TicTacToeGame struct {
//        board         [3][3]rune
//        currentPlayer rune
//        reader        *bufio.Reader
//    }

//    func (game *TicTacToeGame) initialize() {
//        game.board = InitializeBoard()
//        game.currentPlayer = 'X'
//        game.reader = bufio.NewReader(os.Stdin)
//        printWelcomeMessage()
//    }

//    func (game *TicTacToeGame) play() {
//        for {
//            refreshBoard(game.board)
//            game.board, game.currentPlayer, _ = playRound(game.board, game.currentPlayer, game.reader)
//            if ended, message := checkGameEnd(game.board, game.currentPlayer); ended {
//                displayMessage(message)
//                break
//            }
//        }
//    }
//    ```

//    Then, in `main`, you would simply call:
//    ```go
//    func main() {
//        game := &TicTacToeGame{}
//        for {
//            game.initialize()
//            game.play()
//            if !promptRestart(game.reader) {
//                displayMessage("Thank you for playing! Goodbye.")
//                break
//            }
//        }
//    }
//    ```

//    - **Benefit:** Encapsulating the game state and logic into a struct (`TicTacToeGame`) simplifies state management and reduces the complexity in `main`. It also makes it easier to extend the game (e.g., adding a two-player online mode).

// ---

// ### **Summary of Refactorings:**

// 1. **Encapsulate Setup Logic:** Extract game setup into `initializeGame`.
// 2. **Extract Game Loop:** Move the main game loop into `runGameLoop`.
// 3. **Separate Turn Logic:** Split player move input from move execution.
// 4. **Centralize Win/Draw Check:** Use `checkGameEnd` to abstract win and draw conditions.
// 5. **Handle Restart Logic:** Move the restart prompt to `promptRestart`.
// 6. **Abstract Input/Output:** Use helper functions for reading input and displaying output.
// 7. **Encapsulate Game Management:**
