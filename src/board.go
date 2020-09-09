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
const black = "\033[30m"
const colourReset = "\033[0m"

type Board struct {
	board [][]Piece
}

type Piece struct {
	x, y    int
	symbol  string
	isBlack bool
}

type Position struct {
	x, y int
}

func (b *Board) isEmpty(x, y int) bool {
	return b.board[x][y].symbol == " "
}

func (p Piece) checkAllowedMoves(b *Board, newX, newY int) bool {
	currentPiece := b.board[newX][newY]
	if currentPiece.symbol != " " && currentPiece.isBlack == p.isBlack {
		//Square is non empty and has the same colour piece
		return false
	}

	// If moved piece results in check return False
	//Check by checking for check with this piece not in current position
	//do these in this order
	//if in check, check if moving piece fixes check

	allowedMoves := make([]Position, 0)
	fmt.Println("XY", p.x, p.y, newX, newY)
	switch p.symbol {
	case "P":
		// If pawn in start position advance 2
		if p.x == 1 && !p.isBlack && b.isEmpty(p.x+2, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x + 2, p.y})
		} else if p.x == 6 && p.isBlack && b.isEmpty(p.x-2, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x - 2, p.y})
		}

		//Advance 1
		if !p.isBlack && b.isEmpty(p.x+1, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x + 1, p.y})
		} else if b.isEmpty(p.x-1, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x - 1, p.y})
		}

		//Diagonal attack
		if !p.isBlack {
			if p.y-1 >= 0 && !b.isEmpty(p.x+1, p.y-1) && b.board[p.x+1][p.y-1].isBlack {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y - 1})
			} else if p.y+1 <= 7 && !b.isEmpty(p.x+1, p.y+1) && b.board[p.x+1][p.y+1].isBlack {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y + 1})
			}
		} else {
			if p.y-1 >= 0 && !b.isEmpty(p.x-1, p.y-1) && !b.board[p.x-1][p.y-1].isBlack {
				allowedMoves = append(allowedMoves, Position{p.x - 1, p.y - 1})
			} else if p.y+1 <= 7 && !b.isEmpty(p.x-1, p.y+1) && !b.board[p.x-1][p.y+1].isBlack {
				allowedMoves = append(allowedMoves, Position{p.x - 1, p.y + 1})
			}
		}
	case "H":
		//Up
		if p.x+2 <= 7 {
			if p.y+1 <= 7 && (b.isEmpty(p.x+2, p.y+1) || b.board[p.x+2][p.y+1].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 2, p.y + 1})
			}
			if p.y-1 >= 0 && (b.isEmpty(p.x+2, p.y-1) || b.board[p.x+2][p.y-1].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 2, p.y - 1})
			}
		}
		if p.x-2 >= 0 { // Down
			if p.y+1 <= 7 && (b.isEmpty(p.x-2, p.y+1) || b.board[p.x-2][p.y+1].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x - 2, p.y + 1})
			}
			if p.y-1 >= 0 && (b.isEmpty(p.x-2, p.y-1) || b.board[p.x-2][p.y-1].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x - 2, p.y - 1})
			}
		}
		if p.y+2 <= 7 { //Right
			if p.x+1 <= 7 && (b.isEmpty(p.x+1, p.y+2) || b.board[p.x+1][p.y+2].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y + 2})
			}
			if p.x-1 >= 0 && (b.isEmpty(p.x-1, p.y+2) || b.board[p.x-1][p.y+2].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x - 1, p.y + 2})
			}
		}
		if p.y-2 >= 0 { // Left
			if p.x+1 <= 7 && (b.isEmpty(p.x+1, p.y-2) || b.board[p.x+1][p.y-2].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y - 2})
			}
			if p.x-1 >= 0 && (b.isEmpty(p.x-1, p.y-2) || b.board[p.x-1][p.y-2].isBlack != p.isBlack) {
				allowedMoves = append(allowedMoves, Position{p.x - 1, p.y - 2})
			}
		}
	
	case "Q":
		fallthrough
	
	case "R":
		//Up
		counter := 1
		for p.x+counter <= 7 && b.isEmpty(p.x+counter, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x + counter, p.y})
			counter++
		}	
		if p.x+counter <= 7 && b.board[p.x+counter][p.y].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{p.x + counter, p.y})
		}	
		//Down
		counter = 1
		for p.x-counter >= 0 && b.isEmpty(p.x-counter, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x - counter, p.y})
			counter++
		}	
		if p.x-counter >= 0 && b.board[p.x-counter][p.y].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{p.x - counter, p.y})
		}	
		//Right
		counter = 1
		for p.y+counter <= 7 && b.isEmpty(p.x, p.y+counter) {
			allowedMoves = append(allowedMoves, Position{p.x, p.y + counter})
			counter++
		}	
		if p.y+counter <= 7 && b.board[p.x][p.y+counter].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{p.x, p.y + counter})
		}	
		//Left
		counter = 1
		for p.y-counter >= 0 && b.isEmpty(p.x, p.y-counter) {
			allowedMoves = append(allowedMoves, Position{p.x, p.y - counter})
			counter++
		}	
		if p.y-counter >= 0 && b.board[p.x][p.y-counter].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{p.x, p.y - counter})
		}	
		
		if p.symbol == "R" {
			break
		}
		fallthrough
			
	case "B":
		counter := 1
		// Top Right
		tempX, tempY := p.x+counter, p.y+counter;
		for tempX <= 7 && tempY <= 7 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x+counter, p.y+counter;
		}
		if tempX <= 7 && tempY <= 7 && b.board[tempX][tempY].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}

		counter = 1
		// Top Left
		tempX, tempY = p.x+counter, p.y-counter;
		for tempX <= 7 && tempY >= 0 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x+counter, p.y-counter;
		}
		if tempX <= 7 && tempY >= 0 && b.board[tempX][tempY].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}

		counter = 1
		// Bottom Right
		tempX, tempY = p.x-counter, p.y+counter;
		for tempX >= 0 && tempY <= 7 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x-counter, p.y+counter;
		}
		if tempX >= 0 && tempY <= 7 && b.board[tempX][tempY].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}

		counter = 1
		// Bottom Left
		tempX, tempY = p.x-counter, p.y-counter;
		for tempX >= 0 && tempY >= 0 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x-counter, p.y-counter;
		}
		if tempX >= 0 && tempY >= 0 && b.board[tempX][tempY].isBlack != p.isBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}


	case "K":
		

	}
	fmt.Println("moves", allowedMoves)
	for _, val := range allowedMoves {
		if val.x == newX && val.y == newY {
			return true
		}
	}
	fmt.Println("That is an invalid move")
	return false
}

func (p Piece) move(b *Board, newX, newY int) {
	// Check if allowed move
	if p.checkAllowedMoves(b, newX, newY) {
		b.board[p.x][p.y] = Piece{p.x, p.y, " ", false}
		p.x = newX
		p.y = newY
		b.board[newX][newY] = p
		fmt.Println("Moved piece to ", newX, newY, p.x, p.y)
	}

}

//String representation of board
func (b Board) String() string {
	var board string
	var swapColour bool
	board += "\n----------------------------\n"
	// board += "|  a  b  c  d  e  f  g  h  |\n"
	board += "|  0  1  2  3  4  5  6  7  |\n"
	board += "----------------------------\n"
	for row := 7; row >= 0; row-- {
		board += colourReset + strconv.Itoa(row) + "|"
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
		board += colourReset + "|" + strconv.Itoa(row) + "\n"
	}
	board += "----------------------------\n"
	// board += "|  a  b  c  d  e  f  g  h  |\n"
	board += "|  0  1  2  3  4  5  6  7  |\n"
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
			board.board[row][col] = Piece{row, col, " ", false}
		}
	}

	//White pieces
	board.board[0][0] = Piece{0, 0, "R", false}
	board.board[0][1] = Piece{0, 1, "H", false}
	board.board[0][2] = Piece{0, 2, "B", false}
	board.board[0][3] = Piece{0, 3, "Q", false}
	board.board[0][4] = Piece{0, 4, "K", false}
	board.board[0][5] = Piece{0, 5, "B", false}
	board.board[0][6] = Piece{0, 6, "H", false}
	board.board[0][7] = Piece{0, 7, "R", false}

	board.board[1][0] = Piece{1, 0, "P", false}
	board.board[1][1] = Piece{1, 1, "P", false}
	board.board[1][2] = Piece{1, 2, "P", false}
	board.board[1][3] = Piece{1, 3, "P", false}
	board.board[1][4] = Piece{1, 4, "P", false}
	board.board[1][5] = Piece{1, 5, "P", false}
	board.board[1][6] = Piece{1, 6, "P", false}
	board.board[1][7] = Piece{1, 7, "P", false}
	//Black pieces
	board.board[7][0] = Piece{7, 0, "R", true}
	board.board[7][1] = Piece{7, 1, "H", true}
	board.board[7][2] = Piece{7, 2, "B", true}
	board.board[7][3] = Piece{7, 3, "Q", true}
	board.board[7][4] = Piece{7, 4, "K", true}
	board.board[7][5] = Piece{7, 5, "B", true}
	board.board[7][6] = Piece{7, 6, "H", true}
	board.board[7][7] = Piece{7, 7, "R", true}

	board.board[6][0] = Piece{6, 0, "P", true}
	board.board[6][1] = Piece{6, 1, "P", true}
	board.board[6][2] = Piece{6, 2, "P", true}
	board.board[6][3] = Piece{6, 3, "P", true}
	board.board[6][4] = Piece{6, 4, "P", true}
	board.board[6][5] = Piece{6, 5, "P", true}
	board.board[6][6] = Piece{6, 6, "P", true}
	board.board[6][7] = Piece{6, 7, "P", true}
}
