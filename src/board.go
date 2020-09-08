package main

import (
	"fmt"
	"strconv"
)

const boardSize = 8
const bgDark = "\033[48;2;181;136;99m"
const bgLight = "\033[48;2;240;217;181m"
const bgDarkW = "\033[38;5;240;48;2;181;136;99m"
const bgLightW = "\033[38;5;240;48;2;240;217;181m"
const bgDarkB = "\033[30;48;2;181;136;99m"
const bgLightB = "\033[30;48;2;240;217;181m"
const white = "\033[38;5;192m"
const black	 = "\033[30m"
const colourReset = "\033[0m"

type Board struct {
	board [][]Piece
}

//String representation of board
func (b Board) String() string {
	var board string
	var swapColour bool
	board += "\n----------------------------\n"
	board += "|  a  b  c  d  e  f  g  h  |\n"
	board += "----------------------------\n"
	for row := 7; row >= 0; row-- {
		board += colourReset + strconv.Itoa(row+1) + "|"
		for col := 0; col <= 7; col++ {
			piece := b.board[row][col]
			if piece.isBlack {
				if swapColour {
					board += bgDarkB + " " + piece.symbol + " "
				} else {
					board += bgLightB + " " + piece.symbol + " "
				}
			} else {			
				if swapColour {
					board += bgDarkW + " " + piece.symbol + " "
				} else {
					board += bgLightW + " " + piece.symbol + " "
				}
			}
			swapColour = !swapColour
		}
		swapColour = !swapColour
		board += colourReset + "|" + strconv.Itoa(row+1) + "\n"
	}
	board += "----------------------------\n"
	board += "|  a  b  c  d  e  f  g  h  |\n"
	board += "----------------------------\n"
	return fmt.Sprintf("%v", board)
}

func SetupBoard(board *Board) {
	board.board = make([][]Piece, boardSize)
	for i := range board.board {
		board.board[i] = make([]Piece, boardSize)
	}

	for row := 2; row <= 5; row++ {
		for col := 0; col <= 7; col++ {
			board.board[row][col] = Piece{" ", false}
		}
	}

	//White pieces
	board.board[0][0] = Piece{"R", false}
	board.board[0][1] = Piece{"H", false}
	board.board[0][2] = Piece{"B", false}
	board.board[0][3] = Piece{"Q", false}
	board.board[0][4] = Piece{"K", false}
	board.board[0][5] = Piece{"B", false}
	board.board[0][6] = Piece{"H", false}
	board.board[0][7] = Piece{"R", false}

	board.board[1][0] = Piece{"P", false}
	board.board[1][1] = Piece{"P", false}
	board.board[1][2] = Piece{"P", false}
	board.board[1][3] = Piece{"P", false}
	board.board[1][4] = Piece{"P", false}
	board.board[1][5] = Piece{"P", false}
	board.board[1][6] = Piece{"P", false}
	board.board[1][7] = Piece{"P", false}
	//Black pieces
	board.board[7][0] = Piece{"R", true}
	board.board[7][1] = Piece{"H", true}
	board.board[7][2] = Piece{"B", true}
	board.board[7][3] = Piece{"Q", true}
	board.board[7][4] = Piece{"K", true}
	board.board[7][5] = Piece{"B", true}
	board.board[7][6] = Piece{"H", true}
	board.board[7][7] = Piece{"R", true}

	board.board[6][0] = Piece{"P", true}
	board.board[6][1] = Piece{"P", true}
	board.board[6][2] = Piece{"P", true}
	board.board[6][3] = Piece{"P", true}
	board.board[6][4] = Piece{"P", true}
	board.board[6][5] = Piece{"P", true}
	board.board[6][6] = Piece{"P", true}
	board.board[6][7] = Piece{"P", true}

	// for rowVal, row := range board.board {
	// 	for colVal, _ := range row {
	// 		board.board[rowVal][colVal] = Piece{"*", true}
	// 	}
	// }
	fmt.Println(board)
}
