package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type Game struct {
	board Board
	isWhiteTurn bool
}

func (g Game) nextTurn() {
	g.isWhiteTurn = !g.isWhiteTurn
}

//Pre: Takes in a character
//Post: True if letter in A-G
func isLetter(s string) bool {
	return s == "a" || s == "b" || s == "c" || s == "d" || s == "e" || s == "f" || s == "g" || s == "h"
}

//Pre: Takes in a character
//Post: True if number between 0-7
func isNumber(s string) bool {
	return s == "0" || s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7"
}

func (g Game) getTurn(r bufio.Reader) (string, string) {
	if g.isWhiteTurn {
		fmt.Println("White's Turn")
	} else {
		fmt.Println("Black's Turn")
	}
	var str, dst string
	for {
		fmt.Println("Enter the piece you want to move and where to: ")
		str, _ = r.ReadString('\n')
		fmt.Println("String 1", str, len(str))
		str = strings.ToLower(str)
		str = strings.TrimSpace(str)
		fmt.Println("String 2", str, len(str))

		if string(str[2]) == " " {
			dst = str[3:]
			fmt.Println("String 3", dst)
			str = str[:2]
			fmt.Println("String 4", str, len(str))
			dst = strings.TrimSpace(dst)
			fmt.Println("String 5", dst, len(dst))

			if len(dst) == 2 {
				fmt.Println("String 6", dst)
				fmt.Println(string(str), string(dst))
				if isLetter(string(str[0])) && isLetter(string(dst[0])) && isNumber(string(str[1])) && isNumber(string(dst[1])) {
					break
				}
			}
		}
		fmt.Println("That is an invalid move")
		fmt.Println("Please use the format: [LetterNumber LetterNumber] to provide the coordinates of the chosen piece, and where to move it")
	}
	return string(str), string(dst)
}

func main() {
	g := Game{Board{}, true}
	reader := bufio.NewReader(os.Stdin)
	a, b := g.getTurn(*reader)
	fmt.Println(a, b)
	// g.testBoard()
    // fmt.Print("Enter text: ")
    // text, _ := reader.ReadString('\n')
	// fmt.Println(text, len(text))
	// text = strings.TrimSpace(text)
	// fmt.Println(len(text))
	// for i := 0; i < len(text); i++ {
	// 	fmt.Println(string(text[i]))
	// }
}

func (g Game) testBoard() {
	board := g.board 

	fmt.Println("Hello World3!")
	SetupBoard(&board)
	// fmt.Println(board)
	// board.board[1][0].move(&board, 2, 0)
	fmt.Println(board)
	board.board[0][1].move(&board, 2, 2)
	fmt.Println(board)
	// board.board[2][2].move(&board, 4, 3)
	// board.board[7][1].move(&board, 5, 2)
	// board.board[5][2].move(&board, 5, 4)
	// board.board[4][3].move(&board, 5, 4)
	board.board[1][3].move(&board, 3, 3)
	board.board[3][3].move(&board, 4, 3)
	board.board[0][3].move(&board, 2, 3)
	board.board[2][3].move(&board, 4, 5)
	board.board[4][5].move(&board, 6, 4)
	board.board[0][4].move(&board, 1, 3)
	board.board[1][3].move(&board, 3, 3)
	board.board[6][4].move(&board, 5, 4)
	board.board[7][3].move(&board, 5, 5)
	board.board[4][4].move(&board, 3, 4)
	
	fmt.Println("END Is black king currently in check: ", board.kingB.isCheck(&board))
	fmt.Println("END Is white king currently in check: ", board.kingW.isCheck(&board))
	board.board[3][4].move(&board, 2, 4)
	board.board[1][3].move(&board, 2, 4)

	board.board[5][5].move(&board, 5,6)
	// board.board[4][4].move(&board, 3,5)
	// board.board[3][5].move(&board, 2,5)
	// board.board[3][4].move(&board, 3,3)
	// board.board[3][3].move(&board, ,3)
	// board.board[1][3].move(&board, 1,4)
	// board.board[1][4].move(&board, 1,5)
	// board.board[1][5].move(&board, 3,5)
	board.board[7][6].move(&board, 5,5)
	board.board[5][5].move(&board, 3,6)
	board.board[3][6].move(&board, 5,7)
	board.board[5][7].move(&board, 4,5)
	board.board[4][5].move(&board, 3,3)
	board.board[3][3].move(&board, 1,2)
	board.board[1][2].move(&board, 0,4)

	board.board[0][4].move(&board, 1,6)
	board.board[1][6].move(&board, 3,7)

	board.board[1][1].move(&board, 3,1)
	board.board[3][1].move(&board, 4,1)
	fmt.Println("---------------------")
	board.board[6][2].move(&board, 4,2)
	fmt.Println("---------------------")
	board.board[4][1].move(&board, 7,2)
	fmt.Println("---------------------")
	board.board[1][5].move(&board, 3,5)
	board.board[6][7].move(&board, 4,7)
	fmt.Println("---------------------")
	board.board[4][1].move(&board, 7,2)
	board.board[2][4].move(&board, 2,5)
	fmt.Println("---------------------")

	board.board[7][5].move(&board, 5,3)
	fmt.Println("---------------------")

	board.board[7][4].move(&board, 7,6)
	

	fmt.Println("END Is black king currently in check: ", board.kingB.isCheck(&board))
	fmt.Println("END Is white king currently in check: ", board.kingW.isCheck(&board))



	// fmt.Println(board)
	// board.board[6][3].move(&board, 2, 3)
	// board.board[1][2].move(&board, 3, 2)
	fmt.Println(board)

	fmt.Println("---------------------")
	// var fir string 
	// fmt.Scanln(&fir)
	// fmt.Println("This is the string:" , fir, string(fir[4]), len(fir))
	// colorReset := "\033[0m"

    // colorRed := "\033[31m"
    // colorGreen := "\033[32m"
    // colorYellow := "\033[33m"
    // colorBlue := "\033[34m"
    // colorPurple := "\033[35m"
    // colorCyan := "\033[36m"
	// colorWhite := "\033[37m"
	// boardDark := "\033[48;2;181;136;99m"
	// boardLight := "\033[48;2;240;217;181m"
    
    // fmt.Println(boardDark,string(colorRed), "test", string(colorReset))
    // fmt.Println(boardLight, string(colorGreen), "test", string(colorReset))
    // fmt.Println(boardDark, string(colorYellow), "test")
    // fmt.Println(string(colorBlue), "test")
    // fmt.Println(string(colorPurple), "test")
    // fmt.Println(string(colorWhite), "test")
    // fmt.Println(string(colorCyan), "test", string(colorReset))
    // fmt.Println("next")
}