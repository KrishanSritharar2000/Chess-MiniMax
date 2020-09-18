package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type Game struct {
	Board Board
	isWhiteTurn bool
}

func (g *Game) nextTurn() {
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
	return s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" || s == "8"
}

func (g *Game) getTurn(r bufio.Reader) (string, string) {
	var isCheck bool
	if g.isWhiteTurn {
		fmt.Println("White's Turn")
		isCheck = g.Board.kingW.isCheck(&g.Board)
	} else {
		fmt.Println("Black's Turn")
		isCheck = g.Board.kingB.isCheck(&g.Board)
	}
	if isCheck {
		fmt.Println("CHECK")
	}
	var str, dst string
	for {
		fmt.Println("\nEnter the piece you want to move and where to: ")
		str, _ = r.ReadString('\n')
		str = strings.ToLower(str)
		str = strings.TrimSpace(str)

		if string(str[2]) == " " {
			dst = str[3:]
			str = str[:2]
			dst = strings.TrimSpace(dst)

			if len(dst) == 2 {
				if isLetter(string(str[0])) && isLetter(string(dst[0])) && isNumber(string(str[1])) && isNumber(string(dst[1])) {
					break
				} 
			}
		}
		fmt.Println("\nThat is an invalid move")
		fmt.Println("Please use the format: [LetterNumber LetterNumber]")
		fmt.Println("To provide the coordinates [A-G 1-8] of the chosen piece, and where to move it")
	}
	return string(str), string(dst)
}

func (g *Game) makeMove(start, end string) bool {
	var startX, startY, endX, endY int
	startY = int(start[0]) - int('a')
	startX = int(start[1]) - int('1')
	endY = int(end[0]) - int('a')
	endX = int(end[1]) - int('1')

	piece := g.Board.Board[startX][startY]	
	if piece.symbol == " " {
		fmt.Println("There is no piece there!")
		return false
	} else if piece.isBlack == g.isWhiteTurn {
		fmt.Println("That is not your piece, you cannot move it!")
		return false
	}

	result := g.Board.Board[startX][startY].move(&g.Board, endX, endY)
	//moving white piece need to check piece is white
	if result {
		g.nextTurn()
	} else {
		fmt.Println("That move is not allowed")
	}
	return result
}

func StartGame() {
	g := Game{Board{}, true}
	SetupBoard(&g.Board)
	fmt.Println(g.Board)
	reader := bufio.NewReader(os.Stdin)
	for {
		startPos, endPos := g.getTurn(*reader)
		g.makeMove(startPos, endPos)
		fmt.Println(g.Board)
	}
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
	Board := Board{}

	fmt.Println("Hello World3!")
	SetupBoard(&Board)
	// fmt.Println(Board)
	// Board.Board[1][0].move(&Board, 2, 0)
	fmt.Println(Board)
	Board.Board[0][1].move(&Board, 2, 2)
	fmt.Println(Board)
	// Board.Board[2][2].move(&Board, 4, 3)
	// Board.Board[7][1].move(&Board, 5, 2)
	// Board.Board[5][2].move(&Board, 5, 4)
	// Board.Board[4][3].move(&Board, 5, 4)
	Board.Board[1][3].move(&Board, 3, 3)
	Board.Board[3][3].move(&Board, 4, 3)
	Board.Board[0][3].move(&Board, 2, 3)
	Board.Board[2][3].move(&Board, 4, 5)
	Board.Board[4][5].move(&Board, 6, 4)
	Board.Board[0][4].move(&Board, 1, 3)
	Board.Board[1][3].move(&Board, 3, 3)
	Board.Board[6][4].move(&Board, 5, 4)
	Board.Board[7][3].move(&Board, 5, 5)
	Board.Board[4][4].move(&Board, 3, 4)
	
	fmt.Println("END Is black king currently in check: ", Board.kingB.isCheck(&Board))
	fmt.Println("END Is white king currently in check: ", Board.kingW.isCheck(&Board))
	Board.Board[3][4].move(&Board, 2, 4)
	Board.Board[1][3].move(&Board, 2, 4)

	Board.Board[5][5].move(&Board, 5,6)
	// Board.Board[4][4].move(&Board, 3,5)
	// Board.Board[3][5].move(&Board, 2,5)
	// Board.Board[3][4].move(&Board, 3,3)
	// Board.Board[3][3].move(&Board, ,3)
	// Board.Board[1][3].move(&Board, 1,4)
	// Board.Board[1][4].move(&Board, 1,5)
	// Board.Board[1][5].move(&Board, 3,5)
	Board.Board[7][6].move(&Board, 5,5)
	Board.Board[5][5].move(&Board, 3,6)
	Board.Board[3][6].move(&Board, 5,7)
	Board.Board[5][7].move(&Board, 4,5)
	Board.Board[4][5].move(&Board, 3,3)
	Board.Board[3][3].move(&Board, 1,2)
	Board.Board[1][2].move(&Board, 0,4)

	Board.Board[0][4].move(&Board, 1,6)
	Board.Board[1][6].move(&Board, 3,7)

	Board.Board[1][1].move(&Board, 3,1)
	Board.Board[3][1].move(&Board, 4,1)
	fmt.Println("---------------------")
	Board.Board[6][2].move(&Board, 4,2)
	fmt.Println("---------------------")
	Board.Board[4][1].move(&Board, 7,2)
	fmt.Println("---------------------")
	Board.Board[1][5].move(&Board, 3,5)
	Board.Board[6][7].move(&Board, 4,7)
	fmt.Println("---------------------")
	Board.Board[4][1].move(&Board, 7,2)
	Board.Board[2][4].move(&Board, 2,5)
	fmt.Println("---------------------")

	Board.Board[7][5].move(&Board, 5,3)
	fmt.Println("---------------------")

	Board.Board[7][4].move(&Board, 7,6)
	

	fmt.Println("END Is black king currently in check: ", Board.kingB.isCheck(&Board))
	fmt.Println("END Is white king currently in check: ", Board.kingW.isCheck(&Board))



	// fmt.Println(Board)
	// Board.Board[6][3].move(&Board, 2, 3)
	// Board.Board[1][2].move(&Board, 3, 2)
	fmt.Println(Board)

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