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
	Board                        [][]Piece
	kingW, kingB                 Piece
	whiteCheck, blackCheck       bool
	lastPawnMoveW, lastPawnMoveB Position //used for en passant
	castleCheck                  [6]bool  //WhiteKing, BlackKing, WhiteRookLeft, WhiteRookRight, BlackRookLeft, BlackRookRight
}

type Piece struct {
	x, y    int
	Symbol  string
	IsBlack bool
}

type Position struct {
	x, y int
}

func (b *Board) isEmpty(x, y int) bool {
	return b.Board[x][y].Symbol == " "
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsInt(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func (p Piece) isCheck(b *Board) bool {

	//Hoizontal and Vertical
	var horizontal = []string{"R", "Q"}
	counter := 1
	//Up
	for p.x+counter <= 7 && b.isEmpty(p.x+counter, p.y) {
		counter++
	}
	if p.x+counter <= 7 && b.Board[p.x+counter][p.y].IsBlack != p.IsBlack {
		if Contains(horizontal, b.Board[p.x+counter][p.y].Symbol) {
			return true
		}
	}
	//Down
	counter = 1
	for p.x-counter >= 0 && b.isEmpty(p.x-counter, p.y) {
		counter++
	}
	if p.x-counter >= 0 && b.Board[p.x-counter][p.y].IsBlack != p.IsBlack {
		if Contains(horizontal, b.Board[p.x-counter][p.y].Symbol) {
			return true
		}
	}
	//Right
	counter = 1
	for p.y+counter <= 7 && b.isEmpty(p.x, p.y+counter) {
		counter++
	}
	if p.y+counter <= 7 && b.Board[p.x][p.y+counter].IsBlack != p.IsBlack {
		if Contains(horizontal, b.Board[p.x][p.y+counter].Symbol) {
			return true
		}
	}
	//Left
	counter = 1
	for p.y-counter >= 0 && b.isEmpty(p.x, p.y-counter) {
		counter++
	}
	if p.y-counter >= 0 && b.Board[p.x][p.y-counter].IsBlack != p.IsBlack {
		if Contains(horizontal, b.Board[p.x][p.y-counter].Symbol) {
			return true
		}
	}

	//Diagonals
	var diagonal = []string{"B", "Q"}
	counter = 1
	//UpRight
	for p.x+counter <= 7 && p.y+counter <= 7 && b.isEmpty(p.x+counter, p.y+counter) {
		counter++
	}
	if p.x+counter <= 7 && p.y+counter <= 7 && b.Board[p.x+counter][p.y+counter].IsBlack != p.IsBlack {
		if Contains(diagonal, b.Board[p.x+counter][p.y+counter].Symbol) {
			return true
		}
	}
	//DownRight
	counter = 1
	for p.x-counter >= 0 && p.y+counter <= 7 && b.isEmpty(p.x-counter, p.y+counter) {
		counter++
	}
	if p.x-counter >= 0 && p.y+counter <= 7 && b.Board[p.x-counter][p.y+counter].IsBlack != p.IsBlack {
		if Contains(diagonal, b.Board[p.x-counter][p.y+counter].Symbol) {
			return true
		}
	}
	//UpLeft
	counter = 1
	for p.x+counter <= 7 && p.y-counter >= 0 && b.isEmpty(p.x+counter, p.y-counter) {
		counter++
	}
	if p.x+counter <= 7 && p.y-counter >= 0 && b.Board[p.x+counter][p.y-counter].IsBlack != p.IsBlack {
		if Contains(diagonal, b.Board[p.x+counter][p.y-counter].Symbol) {
			return true
		}
	}
	//DownLeft
	counter = 1
	for p.x-counter >= 0 && p.y-counter >= 0 && b.isEmpty(p.x-counter, p.y-counter) {
		counter++
	}
	if p.x-counter >= 0 && p.y-counter >= 0 && b.Board[p.x-counter][p.y-counter].IsBlack != p.IsBlack {
		if Contains(diagonal, b.Board[p.x-counter][p.y-counter].Symbol) {
			return true
		}
	}

	//Checks Pawns
	if p.IsBlack {
		if p.x-1 >= 0 && p.y-1 >= 0 {
			if piece := b.Board[p.x-1][p.y-1]; piece.Symbol == "P" && piece.IsBlack != p.IsBlack {
				return true
			}
		}
		if p.x-1 >= 0 && p.y+1 <= 7 {
			if piece := b.Board[p.x-1][p.y+1]; piece.Symbol == "P" && piece.IsBlack != p.IsBlack {
				return true
			}
		}
	} else {
		if p.x+1 <= 7 && p.y-1 >= 0 {
			if piece := b.Board[p.x+1][p.y-1]; piece.Symbol == "P" && piece.IsBlack != p.IsBlack {
				return true
			}
		}
		if p.x+1 <= 7 && p.y+1 <= 7 {
			if piece := b.Board[p.x+1][p.y+1]; piece.Symbol == "P" && piece.IsBlack != p.IsBlack {
				return true
			}
		}
	}

	//Check horses
	if p.x+2 <= 7 && p.y-1 >= 0 { //Upleft
		if piece := b.Board[p.x+2][p.y-1]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}
	if p.x+2 <= 7 && p.y+1 <= 7 { //UpRight
		if piece := b.Board[p.x+2][p.y+1]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}
	if p.x-2 >= 0 && p.y-1 >= 0 { //DownLeft
		if piece := b.Board[p.x-2][p.y-1]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}
	if p.x-2 >= 0 && p.y+1 <= 7 { //DownRight
		if piece := b.Board[p.x-2][p.y+1]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}
	if p.x+1 <= 7 && p.y+2 <= 7 { //RightUp
		if piece := b.Board[p.x+1][p.y+2]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}
	if p.x-1 >= 0 && p.y+2 <= 7 { //RightDown
		if piece := b.Board[p.x-1][p.y+2]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}
	if p.x+1 <= 7 && p.y-2 >= 0 { //LeftUp
		if piece := b.Board[p.x+1][p.y-2]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}
	if p.x-1 >= 0 && p.y-2 >= 0 { //LeftDown
		if piece := b.Board[p.x-1][p.y-2]; piece.Symbol == "H" && piece.IsBlack != p.IsBlack {
			return true
		}
	}

	//No Check
	return false

}

//Returns true if not check in that place
func (p Piece) checkCastleGap(b *Board, x, y int) bool {
	origX, origY := p.x, p.y
	p.x, p.y = x, y
	if p.IsBlack {
		b.kingB = Piece{x, y, "K", true}
	} else {
		b.kingW = Piece{x, y, "K", false}
	}
	result := p.isCheck(b)
	p.x, p.y = origX, origY
	if p.IsBlack {
		b.kingB = p
	} else {
		b.kingW = p
	}
	return !result
}

func resetBoard(b *Board, p, currentPiece, replacingPiece Piece, newX, newY int) {
	if p.Symbol == "K" {
		if p.IsBlack {
			b.kingB = p
		} else {
			b.kingW = p
		}
	}
	b.Board[p.x][p.y] = currentPiece
	b.Board[newX][newY] = replacingPiece
}

//Pre: newX and newY are in Board boundaries
//Post: Declaration of allowed move
func (p Piece) checkAllowedMoves(b *Board, newX, newY int) bool {
	currentPiece := b.Board[newX][newY]
	if currentPiece.Symbol != " " && currentPiece.IsBlack == p.IsBlack {
		//Square is non empty and has the same colour piece
		fmt.Println("Your piece is there, you cannot move there!")
		return false
	}

	//Checks if moved piece produces check and
	//Checks if moved piece resolves current check
	currentPiece = b.Board[p.x][p.y]
	//Remove the piece
	b.Board[p.x][p.y] = Piece{p.x, p.y, " ", false}
	replaceingPiece := b.Board[newX][newY]
	b.Board[newX][newY] = currentPiece
	if p.Symbol == "K" {
		if p.IsBlack {
			b.kingB = Piece{newX, newY, "K", true}
		} else {
			b.kingW = Piece{newX, newY, "K", false}
		}
	}

	if p.IsBlack {
		if b.kingB.isCheck(b) {
			resetBoard(b, p, currentPiece, replaceingPiece, newX, newY)
			fmt.Println("Invalid Move")
			if p.Symbol == "K" {
				fmt.Println("You will be in check!")
			} else {
				fmt.Println("You are currently in check!")
			}
			return false
		}
	} else if b.kingW.isCheck(b) {
		resetBoard(b, p, currentPiece, replaceingPiece, newX, newY)
		fmt.Println("Invalid Move")
		if p.Symbol == "K" {
			fmt.Println("You will be in check!")
		} else {
			fmt.Println("You are currently in check!")
		}
		return false
	}
	resetBoard(b, p, currentPiece, replaceingPiece, newX, newY)
	// If moved piece results in check return False
	//Check by checking for check with this piece not in current position
	//do these in this order

	//if in check, check if moving piece fixes check

	allowedMoves := p.generatePossibleMoves(b)
	// fmt.Println("XY", p.x, p.y, newX, newY)
	// fmt.Println("WhiteKing", b.kingW, "BlackKing", b.kingB)

	// fmt.Println("moves", allowedMoves)
	// if len(allowedMoves) > 0 {
	// 	fmt.Println("moves First item", allowedMoves[0].x)
	// }
	// fmt.Println("Is black king currently in check: ", b.kingB.isCheck(b))
	// fmt.Println("Is white king currently in check: ", b.kingW.isCheck(b))
	for _, val := range allowedMoves {
		if val.x == newX && val.y == newY {
			return true
		}
	}
	fmt.Println("That is an invalid move")
	return false
}

func (p Piece) generatePossibleMoves(b *Board) []Position {
	allowedMoves := make([]Position, 0)
	// fmt.Println("XY", p.x, p.y, newX, newY)
	// fmt.Println("WhiteKing", b.kingW, "BlackKing", b.kingB)
	fmt.Println("THIS IS THE X AND Y", p.x, p.y)
	switch p.Symbol {
	case "P":
		// If pawn in start position advance 2
		if p.x == 1 && !p.IsBlack && b.isEmpty(p.x+2, p.y) && b.isEmpty(p.x+1, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x + 2, p.y})
		} else if p.x == 6 && p.IsBlack && b.isEmpty(p.x-2, p.y) && b.isEmpty(p.x-1, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x - 2, p.y})
		}

		//Advance 1
		if !p.IsBlack && b.isEmpty(p.x+1, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x + 1, p.y})
		} else if p.IsBlack && b.isEmpty(p.x-1, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x - 1, p.y})
		}

		//Diagonal attack
		if !p.IsBlack {
			if p.y-1 >= 0 && !b.isEmpty(p.x+1, p.y-1) && b.Board[p.x+1][p.y-1].IsBlack {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y - 1})
			} 
			if p.y+1 <= 7 && !b.isEmpty(p.x+1, p.y+1) && b.Board[p.x+1][p.y+1].IsBlack {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y + 1})
			}
		} else {
			if p.y-1 >= 0 && !b.isEmpty(p.x-1, p.y-1) && !b.Board[p.x-1][p.y-1].IsBlack {
				allowedMoves = append(allowedMoves, Position{p.x - 1, p.y - 1})
			} 
			if p.y+1 <= 7 && !b.isEmpty(p.x-1, p.y+1) && !b.Board[p.x-1][p.y+1].IsBlack {
				allowedMoves = append(allowedMoves, Position{p.x - 1, p.y + 1})
			}
		}

		//En passant
		if p.IsBlack {
			if p.x == 3 && (b.lastPawnMoveW.y == p.y-1 || b.lastPawnMoveW.y == p.y+1) {
				allowedMoves = append(allowedMoves, Position{b.lastPawnMoveW.x - 1, b.lastPawnMoveW.y})
			}
		} else {
			// fmt.Println("Last moved black pawn", b.lastPawnMoveB)
			if p.x == 4 && (b.lastPawnMoveB.y == p.y-1 || b.lastPawnMoveB.y == p.y+1) {
				allowedMoves = append(allowedMoves, Position{b.lastPawnMoveB.x + 1, b.lastPawnMoveB.y})
			}
		}
	case "H":
		//Up
		if p.x+2 <= 7 {
			if p.y+1 <= 7 && (b.isEmpty(p.x+2, p.y+1) || b.Board[p.x+2][p.y+1].IsBlack != p.IsBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 2, p.y + 1})
			}
			if p.y-1 >= 0 && (b.isEmpty(p.x+2, p.y-1) || b.Board[p.x+2][p.y-1].IsBlack != p.IsBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 2, p.y - 1})
			}
		}
		if p.x-2 >= 0 { // Down
			if p.y+1 <= 7 && (b.isEmpty(p.x-2, p.y+1) || b.Board[p.x-2][p.y+1].IsBlack != p.IsBlack) {
				allowedMoves = append(allowedMoves, Position{p.x - 2, p.y + 1})
			}
			if p.y-1 >= 0 && (b.isEmpty(p.x-2, p.y-1) || b.Board[p.x-2][p.y-1].IsBlack != p.IsBlack) {
				allowedMoves = append(allowedMoves, Position{p.x - 2, p.y - 1})
			}
		}
		if p.y+2 <= 7 { //Right
			if p.x+1 <= 7 && (b.isEmpty(p.x+1, p.y+2) || b.Board[p.x+1][p.y+2].IsBlack != p.IsBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y + 2})
			}
			if p.x-1 >= 0 && (b.isEmpty(p.x-1, p.y+2) || b.Board[p.x-1][p.y+2].IsBlack != p.IsBlack) {
				allowedMoves = append(allowedMoves, Position{p.x - 1, p.y + 2})
			}
		}
		if p.y-2 >= 0 { // Left
			if p.x+1 <= 7 && (b.isEmpty(p.x+1, p.y-2) || b.Board[p.x+1][p.y-2].IsBlack != p.IsBlack) {
				allowedMoves = append(allowedMoves, Position{p.x + 1, p.y - 2})
			}
			if p.x-1 >= 0 && (b.isEmpty(p.x-1, p.y-2) || b.Board[p.x-1][p.y-2].IsBlack != p.IsBlack) {
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
		if p.x+counter <= 7 && b.Board[p.x+counter][p.y].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{p.x + counter, p.y})
		}
		//Down
		counter = 1
		for p.x-counter >= 0 && b.isEmpty(p.x-counter, p.y) {
			allowedMoves = append(allowedMoves, Position{p.x - counter, p.y})
			counter++
		}
		if p.x-counter >= 0 && b.Board[p.x-counter][p.y].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{p.x - counter, p.y})
		}
		//Right
		counter = 1
		for p.y+counter <= 7 && b.isEmpty(p.x, p.y+counter) {
			allowedMoves = append(allowedMoves, Position{p.x, p.y + counter})
			counter++
		}
		if p.y+counter <= 7 && b.Board[p.x][p.y+counter].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{p.x, p.y + counter})
		}
		//Left
		counter = 1
		for p.y-counter >= 0 && b.isEmpty(p.x, p.y-counter) {
			allowedMoves = append(allowedMoves, Position{p.x, p.y - counter})
			counter++
		}
		if p.y-counter >= 0 && b.Board[p.x][p.y-counter].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{p.x, p.y - counter})
		}

		if p.Symbol == "R" {
			break
		}
		fallthrough

	case "B":
		counter := 1
		// Top Right
		tempX, tempY := p.x+counter, p.y+counter
		for tempX <= 7 && tempY <= 7 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x+counter, p.y+counter
		}
		if tempX <= 7 && tempY <= 7 && b.Board[tempX][tempY].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}

		counter = 1
		// Top Left
		tempX, tempY = p.x+counter, p.y-counter
		for tempX <= 7 && tempY >= 0 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x+counter, p.y-counter
		}
		if tempX <= 7 && tempY >= 0 && b.Board[tempX][tempY].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}

		counter = 1
		// Bottom Right
		tempX, tempY = p.x-counter, p.y+counter
		for tempX >= 0 && tempY <= 7 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x-counter, p.y+counter
		}
		if tempX >= 0 && tempY <= 7 && b.Board[tempX][tempY].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}

		counter = 1
		// Bottom Left
		tempX, tempY = p.x-counter, p.y-counter
		for tempX >= 0 && tempY >= 0 && b.isEmpty(tempX, tempY) {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
			counter++
			tempX, tempY = p.x-counter, p.y-counter
		}
		if tempX >= 0 && tempY >= 0 && b.Board[tempX][tempY].IsBlack != p.IsBlack {
			allowedMoves = append(allowedMoves, Position{tempX, tempY})
		}

	case "K":
		if p.y-1 >= 0 && (b.isEmpty(p.x, p.y-1) || b.Board[p.x][p.y-1].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x, p.y - 1})
		}
		if p.y+1 <= 7 && (b.isEmpty(p.x, p.y+1) || b.Board[p.x][p.y+1].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x, p.y + 1})
		}
		if p.x-1 >= 0 && (b.isEmpty(p.x-1, p.y) || b.Board[p.x-1][p.y].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x - 1, p.y})
		}
		if p.x+1 <= 7 && (b.isEmpty(p.x+1, p.y) || b.Board[p.x+1][p.y].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x + 1, p.y})
		}

		if p.x-1 >= 0 && p.y-1 >= 0 && (b.isEmpty(p.x-1, p.y-1) || b.Board[p.x-1][p.y-1].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x - 1, p.y - 1})
		}
		if p.x+1 <= 7 && p.y+1 <= 7 && (b.isEmpty(p.x+1, p.y+1) || b.Board[p.x+1][p.y+1].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x + 1, p.y + 1})
		}
		if p.x-1 >= 0 && p.y+1 <= 7 && (b.isEmpty(p.x-1, p.y+1) || b.Board[p.x-1][p.y+1].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x - 1, p.y + 1})
		}
		if p.x+1 <= 7 && p.y-1 >= 0 && (b.isEmpty(p.x+1, p.y-1) || b.Board[p.x+1][p.y-1].IsBlack != p.IsBlack) {
			allowedMoves = append(allowedMoves, Position{p.x + 1, p.y - 1})
		}

		//Castling
		if p.IsBlack {
			if !b.castleCheck[1] && !b.castleCheck[5] && b.isEmpty(7, 5) && b.isEmpty(7, 6) {
				if p.checkCastleGap(b, 7, 5) && p.checkCastleGap(b, 7, 6) {
					allowedMoves = append(allowedMoves, Position{7, 6})
				}
			}
			if !b.castleCheck[1] && !b.castleCheck[4] && b.isEmpty(7, 3) && b.isEmpty(7, 2) && b.isEmpty(7, 1) {
				if p.checkCastleGap(b, 7, 3) && p.checkCastleGap(b, 7, 2) && p.checkCastleGap(b, 7, 1) {
					allowedMoves = append(allowedMoves, Position{7, 2})
				}
			}
		} else {
			if !b.castleCheck[0] && !b.castleCheck[3] && b.isEmpty(0, 5) && b.isEmpty(0, 6) {
				if p.checkCastleGap(b, 0, 5) && p.checkCastleGap(b, 0, 6) {
					allowedMoves = append(allowedMoves, Position{0, 6})
				}
			}
			if !b.castleCheck[0] && !b.castleCheck[2] && b.isEmpty(0, 3) && b.isEmpty(0, 2) && b.isEmpty(0, 1) {
				if p.checkCastleGap(b, 0, 3) && p.checkCastleGap(b, 0, 2) && p.checkCastleGap(b, 0, 1) {
					allowedMoves = append(allowedMoves, Position{0, 2})
				}
			}
		}

		//Check if each of the moves dont result in a check
		//Done
		//add castling by checking if king and rooks on starting square and not checks on square in between

	}
	return allowedMoves
}

func (p Piece) removeInvalidMoves(b *Board, slice []Position) []Position {
	//Checks if moved piece produces check and
	//Checks if moved piece resolves current check
	fmt.Println("THIS IS THE MOVES BEFORE INVALID REMOVED", slice)
	var removeIndex []int
	var reset bool
	for i, val := range slice {
		reset = false
		currentPiece := b.Board[p.x][p.y]
		//Remove the piece
		b.Board[p.x][p.y] = Piece{p.x, p.y, " ", false}
		replaceingPiece := b.Board[val.x][val.y]
		b.Board[val.x][val.y] = currentPiece
		if p.Symbol == "K" {
			if p.IsBlack {
				b.kingB = Piece{val.x, val.y, "K", true}
			} else {
				b.kingW = Piece{val.x, val.y, "K", false}
			}
		}

		if p.IsBlack {
			if b.kingB.isCheck(b) {
				fmt.Println("REMOVE THIS", val)
				reset = true
				resetBoard(b, p, currentPiece, replaceingPiece, val.x, val.y)
				removeIndex = append(removeIndex, i)
			}
		} else if b.kingW.isCheck(b) {
			fmt.Println("REMOVE THIS", val)
			reset = true
			resetBoard(b, p, currentPiece, replaceingPiece, val.x, val.y)
			removeIndex = append(removeIndex, i)

		}
		if !reset {
			resetBoard(b, p, currentPiece, replaceingPiece, val.x, val.y)
		}
	}
	if len(removeIndex) > 0 {
		fmt.Println("This is the list of indexs to remove:", removeIndex)
		allowedMoves := make([]Position, 0)
		for i := 0; i < len(slice); i++ {
			if !ContainsInt(removeIndex, i) {
				allowedMoves = append(allowedMoves, slice[i])
			}
		}
		return allowedMoves
	}	
	return slice
}


func (p Piece) move(b *Board, newX, newY int) bool {
	// Check if allowed move
	if p.checkAllowedMoves(b, newX, newY) {
		//Sets pawn positions for en passant
		if p.Symbol == "P" {
			if p.IsBlack && p.x == 6 && newX == 4 {
				b.lastPawnMoveB = Position{newX, newY}
			} else if !p.IsBlack && p.x == 1 && newX == 3 {
				b.lastPawnMoveW = Position{newX, newY}
			}
		} else if p.IsBlack {
			b.lastPawnMoveB = Position{-2, -2}
		} else {
			b.lastPawnMoveW = Position{-2, -2}
		}

		//Sets castle checking bools for moved rooks
		if p.Symbol == "R" {
			if p.IsBlack {
				if !b.castleCheck[4] && p.x == 7 && p.y == 0 {
					b.castleCheck[4] = true
				} else if !b.castleCheck[5] && p.x == 7 && p.y == 7 {
					b.castleCheck[5] = true
				}
			} else {
				if !b.castleCheck[2] && p.x == 0 && p.y == 0 {
					b.castleCheck[2] = true
				} else if !b.castleCheck[3] && p.x == 0 && p.y == 7 {
					b.castleCheck[3] = true
				}
			}
		}
		//Moves king and rook for castling
		if p.Symbol == "K" {
			if p.IsBlack && p.x == 7 && p.y == 4 {
				if newY == 6 {
					b.Board[7][7].x = 7
					b.Board[7][7].y = 5 
					b.Board[7][5] = b.Board[7][7]
					b.Board[7][7] = Piece{7, 7, " ", true}
				} else if newY == 2 {
					b.Board[7][0].x = 7
					b.Board[7][0].y = 3 
					b.Board[7][3] = b.Board[7][0]
					b.Board[7][0] = Piece{7, 0, " ", true}
				}
			} else if !p.IsBlack && p.x == 0 && p.y == 4 {
				if newY == 6 {
					b.Board[0][7].x = 0
					b.Board[0][7].y = 5 
					b.Board[0][5] = b.Board[0][7]
					b.Board[0][7] = Piece{0, 7, " ", false}
				} else if newY == 2 {
					b.Board[0][0].x = 0
					b.Board[0][0].y = 3 
					b.Board[0][3] = b.Board[0][0]
					b.Board[0][0] = Piece{0, 0, " ", false}
				}
			}
		}

		b.Board[p.x][p.y] = Piece{p.x, p.y, " ", false}
		p.x = newX
		p.y = newY
		b.Board[newX][newY] = p

		//Sets king position
		if p.Symbol == "K" {
			if p.IsBlack {
				b.kingB = p
				if !b.castleCheck[1] {
					b.castleCheck[1] = true
				}
			} else {
				b.kingW = p
				if !b.castleCheck[0] {
					b.castleCheck[0] = true
				}
			}
		}

		// fmt.Println("Moved piece to ", newX, newY, p.x, p.y)
		return true
	}
	return false
}

//String representation of Board
func (b Board) String() string {
	var Board string
	var swapColour bool
	Board += "\n----------------------------\n"
	Board += "|  a  b  c  d  e  f  g  h  |\n"
	// Board += "|  0  1  2  3  4  5  6  7  |\n"
	Board += "----------------------------\n"
	for row := 7; row >= 0; row-- {
		Board += colourReset + strconv.Itoa(row+1) + "|"
		for col := 0; col <= 7; col++ {
			piece := b.Board[row][col]
			if piece.IsBlack {
				if swapColour {
					Board += bgDarkB + " " + piece.Symbol + " "
				} else {
					Board += bgLightB + " " + piece.Symbol + " "
				}
			} else {
				if swapColour {
					Board += bgDarkW + " " + piece.Symbol + " "
				} else {
					Board += bgLightW + " " + piece.Symbol + " "
				}
			}
			swapColour = !swapColour
		}
		swapColour = !swapColour
		Board += colourReset + "|" + strconv.Itoa(row+1) + "\n"
	}
	Board += "----------------------------\n"
	Board += "|  a  b  c  d  e  f  g  h  |\n"
	// Board += "|  0  1  2  3  4  5  6  7  |\n"
	Board += "----------------------------\n"
	return fmt.Sprintf("%v", Board)
}

func SetupBoard(Board *Board) {
	Board.Board = make([][]Piece, boardSize)
	for i := range Board.Board {
		Board.Board[i] = make([]Piece, boardSize)
	}

	for row := 2; row <= 5; row++ {
		for col := 0; col <= 7; col++ {
			Board.Board[row][col] = Piece{row, col, " ", false}
		}
	}

	//White pieces
	Board.Board[0][0] = Piece{0, 0, "R", false}
	Board.Board[0][1] = Piece{0, 1, "H", false}
	Board.Board[0][2] = Piece{0, 2, "B", false}
	Board.Board[0][3] = Piece{0, 3, "Q", false}
	Board.Board[0][4] = Piece{0, 4, "K", false}
	Board.Board[0][5] = Piece{0, 5, "B", false}
	Board.Board[0][6] = Piece{0, 6, "H", false}
	Board.Board[0][7] = Piece{0, 7, "R", false}

	Board.Board[1][0] = Piece{1, 0, "P", false}
	Board.Board[1][1] = Piece{1, 1, "P", false}
	Board.Board[1][2] = Piece{1, 2, "P", false}
	Board.Board[1][3] = Piece{1, 3, "P", false}
	Board.Board[1][4] = Piece{1, 4, "P", false}
	Board.Board[1][5] = Piece{1, 5, "P", false}
	Board.Board[1][6] = Piece{1, 6, "P", false}
	Board.Board[1][7] = Piece{1, 7, "P", false}
	//Black pieces
	Board.Board[7][0] = Piece{7, 0, "R", true}
	Board.Board[7][1] = Piece{7, 1, "H", true}
	Board.Board[7][2] = Piece{7, 2, "B", true}
	Board.Board[7][3] = Piece{7, 3, "Q", true}
	Board.Board[7][4] = Piece{7, 4, "K", true}
	Board.Board[7][5] = Piece{7, 5, "B", true}
	Board.Board[7][6] = Piece{7, 6, "H", true}
	Board.Board[7][7] = Piece{7, 7, "R", true}

	Board.Board[6][0] = Piece{6, 0, "P", true}
	Board.Board[6][1] = Piece{6, 1, "P", true}
	Board.Board[6][2] = Piece{6, 2, "P", true}
	Board.Board[6][3] = Piece{6, 3, "P", true}
	Board.Board[6][4] = Piece{6, 4, "P", true}
	Board.Board[6][5] = Piece{6, 5, "P", true}
	Board.Board[6][6] = Piece{6, 6, "P", true}
	Board.Board[6][7] = Piece{6, 7, "P", true}

	Board.kingW = Board.Board[0][4]
	Board.kingB = Board.Board[7][4]
	Board.whiteCheck = false
	Board.blackCheck = false
	// print("Inital Value", Board.whiteCheck)
	// print("Inital Value", Board.lastPawnMoveB.x, Board.lastPawnMoveB.y)
	// print("Inital Value", Board.lastPawnMoveW.x, Board.lastPawnMoveW.y)
}
