package main

import (
	"fmt"
)


func main() {
	board := Board{}

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