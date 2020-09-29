package main

import (
	"fmt"
	// "io/ioutil"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	clients = make(map[string]*User)
	session = make(chan *User)
)

func pairPlayers() {
	var firstPlayer, secondPlayer *User
	for {
		fmt.Println("\n\nStarted loop\n\n")
		firstPlayer = <- session
		fmt.Println("\n\nGot First player", firstPlayer, "\n\n")
		secondPlayer = <- session
		fmt.Println("\n\nGot second player", secondPlayer, "\n\n")
		firstPlayer.whitePlayer = firstPlayer.userID
		firstPlayer.blackPlayer = secondPlayer.userID
		secondPlayer.whitePlayer = firstPlayer.userID
		secondPlayer.blackPlayer = secondPlayer.userID
		secondPlayer.Game = firstPlayer.Game
	}
}

type User struct {
	Game     *Game
	GameMode int //0 local PvP, 1 AI, 2 Online PvP
	userID string// ip address
	whitePlayer, blackPlayer string //ip address of players 
	lastMove string
	undoMoveRequested string
	lastMadeMove, secondLastMove *Move

}

func GetUserAndGame(r *http.Request) (*User, *Game) {
	ip := r.RemoteAddr
	if string(ip[0]) == "[" {
		//parse localHost IP
		// ip = ip[:5] //since localhost is like [::1]:PORT
		ip = "127.0.0.1"
	} else {
		ip = strings.Split(ip, ":")[0] //Splits IPv4 address and port into just address
	}
	if usr, ok := clients[ip]; ok {
		return usr, usr.Game
	}
	clients[ip] = &User{&Game{Board{}, true, &MoveStack{}}, 0, ip, "", "", "", "", &Move{}, &Move{}}
	game := clients[ip].Game
	SetupBoard(&game.Board)
	return clients[ip], game
}

func GetOpponent(u *User) (*User) {
	if (u.GameMode == 2) {
		var opp *User
		if (u.userID == u.whitePlayer) {
			opp, _ = clients[u.blackPlayer]
		} else {
			opp, _ = clients[u.whitePlayer]
		}
		return opp
	}
	log.Print("ERROR requesting opponent")
	return nil
}

func SetGameMode(userIP string, gameMode int) {
	if usr, ok := clients[userIP]; ok {
		usr.GameMode = gameMode
		fmt.Println("\nGAME MODE SET\n")
	}
}

func ClearOnlineGame(user *User) {
	var ok bool
	var opponent *User
	if user.whitePlayer == user.userID {
		opponent, ok = clients[user.blackPlayer]
	} else {
		opponent, ok = clients[user.whitePlayer]
	}
	if ok {
		user.whitePlayer = ""
		user.blackPlayer = ""
		opponent.whitePlayer = ""
		opponent.blackPlayer = ""
		fmt.Println("Cleared game opponent data")
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	usr, _ := GetUserAndGame(r)
	
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("website/index.html")
		if err != nil {
			log.Print("Error parsing template: ", err)
		}
		err = t.Execute(w, usr)
		if err != nil {
			log.Print("Error during executing: ", err)
		}
		// if usr.whitePlayer != "" {
		// 	ClearOnlineGame(usr)
		// }
	case "POST":
		r.ParseMultipartForm(0)
		message := r.FormValue("option")
		fmt.Println("Client in HomePage POST:", message)
		num, err := strconv.Atoi(message)
		if err != nil {
			log.Print("Error with message from Client: ", err)
		}
		if num == 2 {
			fmt.Println("Find Opponent")
			SetGameMode(usr.userID, num)
			session <- usr
			for usr.blackPlayer == "" {}
			fmt.Println("These are the users:", usr)
		} else {
			SetGameMode(usr.userID, num)
		}
		fmt.Fprintf(w, "success")
		fmt.Println("To Client: success GameMode:", num, "User:", usr)
	}
}

// Message will be either:
// "opt 00" for give move options for piece at (0,0)
// "mov 00 01" for move piece at (0,0) to (0,1)
// "pwn 00q" to promote pawn at 00 to a queen
// "rst " to restart the game
// "ply " to get player turn and check information
// "bck " to undo turn
// "col" to get player colour and game mode information 
// "opp" to get opponent move
// "aim" to get AI move
func GamePage(w http.ResponseWriter, r *http.Request) {
	usr, game := GetUserAndGame(r)
	fmt.Println("IP:", r.RemoteAddr)
	fmt.Println(clients)
	fmt.Println("GameMode:", usr.GameMode)

	switch r.Method {
	case "GET":
		fmt.Println("GET REQUEST")
		tmpl, err := template.New("game.html").Funcs(template.FuncMap{
			"minus": func(a, b int) int {
				return a - b
			},
			"add": func(a, b int) int {
				return a + b
			},
		}).ParseFiles("website/game.html")
		if err != nil {
			log.Print("Error parsing template: ", err)
		}
		err = tmpl.Execute(w, usr)
		if err != nil {
			log.Print("Error during executing: ", err)
		}
	case "POST":
		r.ParseMultipartForm(0)
		message := r.FormValue("white")
		if len(message) == 0 {
			message = r.FormValue("black")
		}
		if len(message) == 0 {
			message = r.FormValue("empty")
		}
		fmt.Println("----------------------------------")
		fmt.Println("Message from Client:", message, len(message))
		opcode := message
		rest := ""
		if len(message) > 3 {
			splitIndex := strings.Index(message, " ")
			if splitIndex == -1 {
				log.Print("Error in HTTP request:", message, len(message))
			}
			fmt.Println("message:", message, " splitIndex:", splitIndex)
			rest = message[(splitIndex + 1):]
			opcode = message[:splitIndex]
		}

		switch opcode {
		case "opt":
			x, y := int(rest[0])-int('0'), int(rest[1])-int('0')
			fmt.Println("Give the move options for the piece at:", x, y, game.Board.Board[x][y])
			possibleMoves := game.Board.Board[x][y].generatePossibleMoves(&game.Board)
			possibleMoves = game.Board.Board[x][y].removeInvalidMoves(&game.Board, possibleMoves)
			fmt.Println("These are the possible moves:\n", possibleMoves)
			fmt.Fprintf(w, getString(possibleMoves))
		case "mov":
			startX, startY := int(rest[0])-int('0'), int(rest[1])-int('0')
			destX, destY := int(rest[2])-int('0'), int(rest[3])-int('0')
			enPassantCheckPiece := game.Board.Board[destX][destY].Symbol
			fmt.Println("Move the piece at:", startX, startY, game.Board.Board[startX][startY], "to:", destX, destY, game.Board.Board[destX][destY])
			result := game.makeMove(startX, startY, destX, destY)
			fmt.Println("Is there a pawn promotion: ", game.Board.promotePawn)

			checkText := getCheckMessage(game, result)
			//check for en passant
			if result && game.Board.Board[destX][destY].Symbol == "P" && abs(startY-destY) == 1 && enPassantCheckPiece == " " {
				fmt.Fprintf(w, "Result:"+"enpassant"+strconv.Itoa(startX)+strconv.Itoa(destY)+checkText)
				usr.lastMove = "Result:"+"enpassant"+strconv.Itoa(startX)+strconv.Itoa(destY)+checkText
			} else if result && game.Board.Board[destX][destY].Symbol == "K" && abs(startY-destY) == 2 {
				var rookLocation string
				if game.Board.Board[destX][destY].IsBlack {
					if destY == 6 {
						rookLocation = "7775"
					} else if destY == 2 {
						rookLocation = "7073"
					}
				} else {
					if destY == 6 {
						rookLocation = "0705"
					} else if destY == 2 {
						rookLocation = "0003"
					}
				}
				fmt.Fprintf(w, "Result:"+"castle"+rookLocation+checkText)
				usr.lastMove = "Result:"+"castle"+rookLocation+checkText
			} else if game.Board.promotePawn && usr.GameMode != 1 && usr.Game.IsWhiteTurn == true {
				game.Board.promotePawn = false
				fmt.Println("CHANGED PAWN PROMOTE BACK TO: ", game.Board.promotePawn)
				fmt.Fprintf(w, "Result:"+"pwn"+checkText)
				usr.lastMove = "Result:"+"pwn"+checkText
			} else {
				fmt.Println("Result:" + strconv.FormatBool(result) + checkText)
				fmt.Fprintf(w, "Result:"+strconv.FormatBool(result)+checkText)
				usr.lastMove = "Result:"+strconv.FormatBool(result)+checkText
			}
			if game.Board.promotePawn {
				game.Board.promotePawn = false
			}

			fmt.Println("This is the move stack:", game.Moves)
		case "pwn":
			x, y := int(rest[0])-int('0'), int(rest[1])-int('0')
			fmt.Println("This is the string of the oanw to promote:", string(rest[2]))
			result := game.Board.Board[x][y].promotePawn(&game.Board, string(rest[2]))
			checkText := getCheckMessage(game, result)
			fmt.Fprintf(w, strconv.FormatBool(result)+checkText)
		case "rst":
			SetupBoard(&game.Board)
			game.IsWhiteTurn = true
			game.Moves = &MoveStack{}
			fmt.Fprintf(w, "reload")
		case "ply":
			fmt.Fprintf(w, strconv.FormatBool(game.IsWhiteTurn)+strconv.Itoa(usr.Game.Moves.Size())+getCheckMessage(game, true))
		case "bck":
			if usr.GameMode == 2 {
				fmt.Println("Requested undo move from:", r.RemoteAddr)
				usr.undoMoveRequested = "true"
				opp := GetOpponent(usr)
				for opp.undoMoveRequested == "" {}
				opp.undoMoveRequested = ""
				if usr.undoMoveRequested == "rtb" {
					fmt.Fprintf(w, "reject")
					usr.undoMoveRequested = ""
					return
				}
				usr.undoMoveRequested = ""
				fmt.Println(r.RemoteAddr, "take back accepted by opponent")
				move := usr.lastMadeMove
				move2 := usr.secondLastMove
				fmt.Fprintf(w, "true" + strconv.Itoa(move.From.x) + strconv.Itoa(move.From.y) + string(move.From.Symbol) + string(strconv.FormatBool(move.From.IsBlack)[0]) + strconv.Itoa(move.To.x) + strconv.Itoa(move.To.y) + string(move.To.Symbol) + string(strconv.FormatBool(move.To.IsBlack)[0]) + strconv.Itoa(move2.From.x) + strconv.Itoa(move2.From.y) + string(move2.From.Symbol) + string(strconv.FormatBool(move2.From.IsBlack)[0]) + strconv.Itoa(move2.To.x) + strconv.Itoa(move2.To.y) + string(move2.To.Symbol) + string(strconv.FormatBool(move2.To.IsBlack)[0]) + "ac")
				usr.lastMadeMove = &Move{}
				usr.secondLastMove = &Move{}
				return
			}
			result, move := game.undoTurn()
			if result {
				fmt.Println("Un did this move:", move)
				if usr.GameMode == 1 {
					result2, move2 := game.undoTurn()
					if result2 { 
						fmt.Fprintf(w, "true" + strconv.Itoa(move.From.x) + strconv.Itoa(move.From.y) + string(move.From.Symbol) + string(strconv.FormatBool(move.From.IsBlack)[0]) + strconv.Itoa(move.To.x) + strconv.Itoa(move.To.y) + string(move.To.Symbol) + string(strconv.FormatBool(move.To.IsBlack)[0]) + strconv.Itoa(move2.From.x) + strconv.Itoa(move2.From.y) + string(move2.From.Symbol) + string(strconv.FormatBool(move2.From.IsBlack)[0]) + strconv.Itoa(move2.To.x) + strconv.Itoa(move2.To.y) + string(move2.To.Symbol) + string(strconv.FormatBool(move2.To.IsBlack)[0]) + "ai")
					} else {
						fmt.Fprintf(w, "false")
					}
				} else {
					fmt.Fprintf(w, "true" + strconv.Itoa(move.From.x) + strconv.Itoa(move.From.y) + string(move.From.Symbol) + string(strconv.FormatBool(move.From.IsBlack)[0]) + strconv.Itoa(move.To.x) + strconv.Itoa(move.To.y) + string(move.To.Symbol) + string(strconv.FormatBool(move.To.IsBlack)[0]))
				}
				
			} else {
				fmt.Fprintf(w, "false")
			}
		case "col":
			fmt.Println("User:", usr)
			fmt.Println("Gamemode:", usr.GameMode)
			if usr.GameMode == 2 {
				fmt.Println("User:", usr)
				fmt.Println("response:",strconv.FormatBool(usr.whitePlayer == usr.userID))
				fmt.Fprintf(w, strconv.FormatBool(usr.whitePlayer == usr.userID) + strconv.Itoa(usr.GameMode))
			} else if usr.GameMode == 1 {
				//change later
				fmt.Fprintf(w, "true" + strconv.Itoa(usr.GameMode))	
			} else {
				fmt.Fprintf(w, "true" + strconv.Itoa(usr.GameMode))
			}
		case "opp":
			// for {
			// 	val, ok := usr.Game.Moves.Peek(); 
			// 	if ok && val.From.Symbol != " " && val.From.IsBlack == !game.IsWhiteTurn {
			// 		return
			// 	}
			// }
			fmt.Println("GOT OPP REQUEST FROM", r.RemoteAddr)
			fmt.Println("Current Turn:", game.IsWhiteTurn, usr.userID == usr.whitePlayer, usr.userID == usr.blackPlayer)
			// for game.IsWhiteTurn != (usr.userID == usr.whitePlayer) {}
			opp := GetOpponent(usr)
			for (opp.lastMove == "" || strings.Contains(opp.lastMove, "false")) && opp.undoMoveRequested == "" {}
			if opp.undoMoveRequested == "" {
				lastMove := opp.lastMove
				fmt.Println("This is the last move:", lastMove)
				fmt.Println("This is the user:", usr)
				fmt.Println("This is the opp:", opp)
				opp.lastMove = ""
				val, ok := usr.Game.Moves.Peek()
				fmt.Println("This is the val and ok:", val, ok)
				if ok {
					fmt.Fprintf(w, lastMove + strconv.Itoa(val.From.x) + strconv.Itoa(val.From.y) + strconv.Itoa(val.To.x) + strconv.Itoa(val.To.y))
				} else {
					fmt.Fprintf(w, "false")
				}
				fmt.Println("This is the last move:", lastMove + strconv.Itoa(val.From.x) + strconv.Itoa(val.From.y) + strconv.Itoa(val.To.x) + strconv.Itoa(val.To.y))
			} else {
				fmt.Println("Asked ", r.RemoteAddr, "for modal")
				fmt.Fprintf(w, "bck")
			}
		case "aim":
			fmt.Println("GOT AIM REQUEST FROM", r.RemoteAddr)
			start := time.Now()
			moveCounter := 0
			for usr.Game.IsWhiteTurn == true {}
			AIMove := usr.Game.FindBestMove(false, 3, &moveCounter)
			usr.Game.makeMove(AIMove.From.x, AIMove.From.y, AIMove.To.x, AIMove.To.y)
			fmt.Println("Time Taken:", time.Since(start), "seconds")
			fmt.Fprintf(w, "Result:true" + getCheckMessage(usr.Game, true) + strconv.Itoa(AIMove.From.x) + strconv.Itoa(AIMove.From.y) + strconv.Itoa(AIMove.To.x) + strconv.Itoa(AIMove.To.y))
		case "atb":
			fmt.Println(r.RemoteAddr + "accepted")
			opp := GetOpponent(usr)
			//Assumes there are at least two moves on the move stack
			_, move := game.undoTurn()
			_, move2 := game.undoTurn()
			opp.lastMadeMove = &move
			opp.secondLastMove = &move2
			fmt.Println("These are the moves", move, move2)
			opp.undoMoveRequested = "atb"
			usr.undoMoveRequested = "answered"
			fmt.Fprintf(w, "true" + strconv.Itoa(move.From.x) + strconv.Itoa(move.From.y) + string(move.From.Symbol) + string(strconv.FormatBool(move.From.IsBlack)[0]) + strconv.Itoa(move.To.x) + strconv.Itoa(move.To.y) + string(move.To.Symbol) + string(strconv.FormatBool(move.To.IsBlack)[0]) + strconv.Itoa(move2.From.x) + strconv.Itoa(move2.From.y) + string(move2.From.Symbol) + string(strconv.FormatBool(move2.From.IsBlack)[0]) + strconv.Itoa(move2.To.x) + strconv.Itoa(move2.To.y) + string(move2.To.Symbol) + string(strconv.FormatBool(move2.To.IsBlack)[0]) + "ac")
			undoMove(w, usr, game, "ac")
		case "rtb":
			fmt.Println(r.RemoteAddr + "rejected")
			opp := GetOpponent(usr)
			opp.undoMoveRequested = "rtb"
			usr.undoMoveRequested = "answered"
			fmt.Fprintf(w, "reject")
		default:
			log.Print("HTTP Request Error")
		}

		// respond to client's request
		// fmt.Fprintf(w, "Server: %s \n", message+" | "+time.Now().Format(time.RFC3339))
	}
}

func undoMove(w http.ResponseWriter, usr *User, game *Game, gameMode2acceptedText string) {

}

func getCheckMessage(game *Game, result bool) string {
	var kingInCheck, kingInCheckMate, kingInStaleMate bool
	if game.IsWhiteTurn {
		kingInCheck = game.Board.kingW.isCheck(&game.Board)
	} else {
		kingInCheck = game.Board.kingB.isCheck(&game.Board)
	}
	if kingInCheck {
		if game.IsWhiteTurn {
			kingInCheckMate = game.Board.kingW.isCheckMate(&game.Board)
		} else {
			kingInCheckMate = game.Board.kingB.isCheckMate(&game.Board)
		}
	} else {
		if game.IsWhiteTurn {
			kingInStaleMate = game.Board.kingW.isCheckMate(&game.Board)
		} else {
			kingInStaleMate = game.Board.kingB.isCheckMate(&game.Board)
		}
	}

	checkText := ""
	fmt.Println("Checkmate:", kingInCheckMate)
	if result && kingInCheckMate {
		checkText = "mate"
	} else if result && kingInCheck {
		checkText = "check"
	} else if result && kingInStaleMate {
		checkText = "stale"
	}
	return checkText
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func getString(slice []Position) string {
	word := ""
	fmt.Println("Start")
	for _, val := range slice {
		word += strconv.Itoa(val.x) + strconv.Itoa(val.y)
	}
	fmt.Println("this is the word from getString", word)
	return word
}

func main() {
	go pairPlayers()
	// StartGame()
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./website"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/game", GamePage)

	log.Fatal(http.ListenAndServe(":8080", nil))
	// res, err := http.Get("http://www.google.com/robots.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// robots, err := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("---------------------")
	// fmt.Printf("%s", robots)
	// fmt.Println("---------------------")

}
