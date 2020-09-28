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

func (g *Game) FindBestMove(forWhite bool, depth int, compCounter *int) MovePair {
	bestMoves := make([]MovePair, 0)
	var currValue, bestValue int
	if forWhite {
		bestValue = -maxScore
	} else {
		bestValue = maxScore
	}
	moves := g.GetAvailableMoves(forWhite)
	for _, move := range moves {
		if g.makeMove(move.From.x, move.From.y, move.To.x, move.To.y) {
			currValue = g.Minimax(!forWhite, 0, depth, -1000000, 1000000, compCounter)
			g.undoTurn()
		} else {
			fmt.Println(g.Board)
			fmt.Println("Turn:", g.IsWhiteTurn, move.From.x, move.From.y, move.To.x, move.To.y)

		}
		if (forWhite && currValue >= bestValue) || (!forWhite && currValue <= bestValue) {
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
	fmt.Println("Is it whites turn:", g.IsWhiteTurn)
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
func (g *Game) Minimax(isMaxTurn bool, depth, maxDepth, alpha, beta int, compCounter *int) int {
	value := g.GetValue(isMaxTurn)
	*compCounter++

	if value == 1 {
		if isMaxTurn {
			return maxScore - depth
		}
		return depth - maxScore
	} else if value == 2 {
		if isMaxTurn {
			return ((3*maxScore)/4) - depth
		}
		return depth - ((3*maxScore)/4)
	}

	if depth == maxDepth {
		// if isMaxTurn {
		// 	return depth - g.getValueOfPiecesOnBoard(isMaxTurn)
		// }
		// return g.getValueOfPiecesOnBoard(!isMaxTurn) - depth
		if isMaxTurn {
			return g.getValueOfPiecesOnBoard(isMaxTurn) - g.getValueOfPiecesOnBoard(!isMaxTurn) + depth
		}
		return (g.getValueOfPiecesOnBoard(!isMaxTurn) - g.getValueOfPiecesOnBoard(isMaxTurn)) - depth
	}

	var currValue int
	if isMaxTurn {
		best := -1000000
		currValue = -maxScore
		moves := g.GetAvailableMoves(true)
		for _, move := range moves {
			//Make move
			if g.makeMove(move.From.x, move.From.y, move.To.x, move.To.y) {
				//Recurse
				currValue = g.Minimax(false, depth + 1, maxDepth, alpha, beta, compCounter)
				//Undo move
				g.undoTurn()
				best = Max(best, currValue)
				alpha = Max(alpha, currValue) 
				if beta <= alpha {
					break
				}
			} else {
				fmt.Println(g.Board)
				fmt.Println(moves)
				fmt.Println("Turn:", g.IsWhiteTurn, move.From.x, move.From.y, move.To.x, move.To.y)
			}					
		}
	} else {
		best := 1000000
		currValue = maxScore
		moves := g.GetAvailableMoves(false)
		for _, move := range moves {
			//Make move
			if g.makeMove(move.From.x, move.From.y, move.To.x, move.To.y) {
				//Recurse
				currValue = g.Minimax(true, depth + 1, maxDepth, alpha, beta, compCounter)
				//Undo move
				g.undoTurn()
				best = Min(best, currValue)
				beta = Min(beta, currValue)
				if beta <= alpha {
					break
				}

			} else {	
				fmt.Println(g.Board)
				fmt.Println(moves)
				fmt.Println("Turn:", g.IsWhiteTurn, move.From.x, move.From.y, move.To.x, move.To.y)
			}
		}
	}
	return currValue  
}

type quad struct {
	a,b,c,d int
}

func main() {
	g := Game{Board{}, true, &MoveStack{}}
	SetupBoard(&g.Board)
	g.GetAvailableMoves(true)
	fmt.Println("Value:", g.GetValue(true))
	fmt.Println(g.Board)
	reader := bufio.NewReader(os.Stdin)
	// ms := []quad{quad{1,3,3,3},quad{1,4,3,4},quad{1,0,2,0},quad{1,7,2,7}}
	counter := 0
	for {
		for {
			g.GetAvailableMoves(false)
			fmt.Println(g.Board)
			result := g.makeMove(g.getTurn(*reader))
			// ms[counter].a,ms[counter].b,ms[counter].c,ms[counter].d)
			if result {
				break
			}
		}
		start := time.Now()
		moveCounter := 0
		AIMove := g.FindBestMove(false, 3, &moveCounter)
		fmt.Println("This is the AI Move:", AIMove)
		fmt.Println("This is how many comparisons it took:", moveCounter)
		if !g.makeMove(AIMove.From.x, AIMove.From.y, AIMove.To.x, AIMove.To.y) {
			fmt.Println(g.Board)
			fmt.Println(AIMove.From.x, AIMove.From.y, AIMove.To.x, AIMove.To.y)
			return
		}
		fmt.Println(g.Board)
		fmt.Println("Time Taken:", time.Since(start), "seconds")
		counter++
	}
}