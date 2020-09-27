package main

import (
	"fmt"
)

const maxScore = 100
var (
	pieceScores = map[string]int{"Q":9, "R":5, "B":3, "H":3, "P":1}
)

type MovePair struct {
	From, To Position
}

//  Returns:
// 1 for checkmate
// 2 for check
// 0 for nothing
func (g *Game) GetValue(movedWhite bool) int {
	
	//check
	if movedWhite {
		if (g.Board.kingB.isCheck(&g.Board)) {
			//checkmate
			if (g.Board.kingB.isCheckMate(&g.Board)) {
				return 1
			}
			return 2
		}
	} else {
		if (g.Board.kingW.isCheck(&g.Board)) {
			if (g.Board.kingW.isCheckMate(&g.Board)) {
				return 1
			}
			return 2
		}
	}
	return 0
}

func (g *Game) getValueOfPiecesOnBoard(forWhite bool) int {
	score := 0
	pieces := g.GetAllPieces(forWhite)
	for _, piece := range pieces {
		if piece.Symbol != " " && piece.Symbol != "K" {
			pieceScore, _ := pieceScores[piece.Symbol]
			score += pieceScore
		}
	}
	return score
}

func (g *Game) GetAllPieces(forWhite bool) []Piece {
	pieces := make([]Piece, 0)
	for _, row := range g.Board.Board {
		for _, piece := range row {
			if piece.Symbol != " " && piece.IsBlack == !forWhite {
				pieces = append(pieces, piece)
			}
		}
	}
	return pieces
}

func (g *Game) GetAvailableMoves(forWhite bool) []MovePair {
	pieces := g.GetAllPieces(forWhite)
	moves := make([]MovePair, 0)
	for _, piece := range pieces {
		positions := piece.removeInvalidMoves(&g.Board, piece.generatePossibleMoves(&g.Board))
		for _, position := range positions {
			moves = append(moves, MovePair{Position{piece.x, piece.y}, position})
		}
	}
	fmt.Println("These are the available moves: ", len(moves), moves)
	return moves
}

func (g *Game) FindBestMove() {
	fmt.Println("FOUND IT")
}

//White is the maximising player
func (g *Game) Minimax(depth, maxDepth int, isMaxTurn bool) int {
	fmt.Println("In Minimax")
	value := g.GetValue(isMaxTurn)

	if value == 1 {
		if isMaxTurn {
			return depth - maxScore
		}
		return maxScore - depth
	} else if value == 2 {
		if isMaxTurn {
			return depth - ((3*maxScore)/4)
		}
		return ((3*maxScore)/4) - depth
	}

	if depth == maxDepth {
		return g.getValueOfPiecesOnBoard(isMaxTurn)
	}

	if isMaxTurn {
		// best = -1000000
		moves := g.GetAvailableMoves()

	} else {
		// best = 1000000

	}
}

func main() {
	g := Game{Board{}, true, &MoveStack{}}
	SetupBoard(&g.Board)
	g.GetAvailableMoves(true)
	fmt.Println("Value:", g.GetValue(true))
}