package main

import (
	"fmt"
	"math/rand"
	"time"
	"bufio"
	"os"
)

const maxScore = 100
var (
	pieceScores = map[string]int{"Q":9, "R":5, "B":3, "H":3, "P":1}
	randomNumbers = rand.New(rand.NewSource(time.Now().UnixNano()))
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

func (g *Game) FindBestMove(forWhite bool, depth int) MovePair {
	bestMoves := make([]MovePair, 0)
	var currValue, bestValue int
	if forWhite {
		bestValue = maxScore
	} else {
		bestValue = -maxScore
	}
	moves := g.GetAvailableMoves(forWhite)
	for _, move := range moves {
		g.makeMove(move.From.x, move.From.y, move.To.x, move.To.y)
		currValue = g.Minimax(0, depth, !forWhite)
		g.undoTurn()
		if (forWhite && currValue <= bestValue) || (!forWhite && currValue >= bestValue) {
			if currValue == bestValue {
				bestMoves = append(bestMoves, move)
			} else {
				bestValue = currValue
				bestMoves = nil
				bestMoves = make([]MovePair, 0)
				bestMoves = append(bestMoves, move)
			}
			
		}
	}
	fmt.Println("This is the bestValue:", bestValue)
	fmt.Println("These are the bestMoves:", bestMoves)
	if len(bestMoves) > 1 {
		return bestMoves[randomNumbers.Intn(len(bestMoves))]
	}
	return bestMoves[0]

}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	var currValue int
	if isMaxTurn {
		// best = -1000000
		moves := g.GetAvailableMoves(true)
		for _, move := range moves {
			//Make move
			g.makeMove(move.From.x, move.From.y, move.To.x, move.To.y)
			//Recurse
			currValue = Max(currValue, g.Minimax(depth + 1, maxDepth, false))
			//Undo move
			g.undoTurn()
		}
	} else {
		// best = 1000000
		moves := g.GetAvailableMoves(false)
		for _, move := range moves {
			//Make move
			g.makeMove(move.From.x, move.From.y, move.To.x, move.To.y)
			//Recurse
			currValue = Min(currValue, g.Minimax(depth + 1, maxDepth, true))
			//Undo move
			g.undoTurn()
		}
	}
	return currValue  
}

func main() {
	g := Game{Board{}, true, &MoveStack{}}
	SetupBoard(&g.Board)
	g.GetAvailableMoves(true)
	fmt.Println("Value:", g.GetValue(true))
	fmt.Println(g.Board)
	reader := bufio.NewReader(os.Stdin)
	for {
		for {
			fmt.Println(g.Board)
			result := g.makeMove(g.getTurn(*reader))
			if result {
				break
			}
		}
		AIMove := g.FindBestMove(false, 2)
		fmt.Println("This is the AI Move:", AIMove)
		g.makeMove(AIMove.From.x, AIMove.From.y, AIMove.To.x, AIMove.To.y)
		fmt.Println(g.Board)
	}
}