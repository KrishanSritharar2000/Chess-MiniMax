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
	"os"
)

var (
	clients = make(map[string]*User)
	session = make(chan *User)
)

func pairPlayers() {
	var firstPlayer, secondPlayer *User
	for {
		fmt.Println("\n\nStarted loop\n\n")
		firstPlayer = <-session
		fmt.Println("\n\nGot First player", firstPlayer, "\n\n")
		secondPlayer = <-session
		fmt.Println("\n\nGot second player", secondPlayer, "\n\n")
		if firstPlayer.UserID == secondPlayer.UserID {
			//cannot play against yourself
			// session <- firstPlayer
			fmt.Println("same player tried to pair up")
			continue
		}
		if (firstPlayer.WhitePlayer == "" && secondPlayer.WhitePlayer == "") ||
			(firstPlayer.BlackPlayer == "" && secondPlayer.BlackPlayer == "") {
			//chosen to be the same colour
			if firstPlayer.UserID == firstPlayer.WhitePlayer {
				firstPlayer.BlackPlayer = secondPlayer.UserID
				secondPlayer.WhitePlayer = firstPlayer.UserID
				secondPlayer.BlackPlayer = secondPlayer.UserID
			} else {
				firstPlayer.WhitePlayer = secondPlayer.UserID
				secondPlayer.BlackPlayer = firstPlayer.UserID
				secondPlayer.WhitePlayer = secondPlayer.UserID
			}
		} else {
			if firstPlayer.UserID == firstPlayer.WhitePlayer {
				firstPlayer.BlackPlayer = secondPlayer.UserID
				secondPlayer.WhitePlayer = firstPlayer.UserID
			} else {
				firstPlayer.WhitePlayer = secondPlayer.UserID
				secondPlayer.BlackPlayer = firstPlayer.UserID
			}
		}
		firstPlayer.Game = secondPlayer.Game
		fmt.Println("assgined varaibles")
		firstPlayer.findOpponent <- secondPlayer
		secondPlayer.findOpponent <- firstPlayer
	}
}

type User struct {
	Game                     *Game
	GameMode                 int    //0 local PvP, 1 AI, 2 Online PvP
	UserID                   string // ip address
	WhitePlayer, BlackPlayer string //ip address of players
	// undoMoveRequested string
	// lastMadeMove, secondLastMove *Move
	findOpponent  chan *User
	opponenetMove chan string
	undoMove      chan *Move
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
	clients[ip] = &User{&Game{Board{}, true, &MoveStack{}}, 0, ip, "", "", make(chan *User), make(chan string), make(chan *Move)}
	game := clients[ip].Game
	SetupBoard(&game.Board)
	return clients[ip], game
}

func GetOpponent(u *User) *User {
	if u.GameMode == 2 {
		var opp *User
		if u.UserID == u.WhitePlayer {
			opp, _ = clients[u.BlackPlayer]
		} else {
			opp, _ = clients[u.WhitePlayer]
		}
		return opp
	}
	log.Print("ERROR requesting opponent")
	return nil
}

func SetGameModeAndColour(userIP string, gameMode int, isWhite bool) {
	fmt.Println("Chose to be white:", isWhite)
	if usr, ok := clients[userIP]; ok {
		usr.Game = &Game{Board{}, true, &MoveStack{}}
		SetupBoard(&usr.Game.Board)
		usr.GameMode = gameMode
		if isWhite {
			usr.WhitePlayer = usr.UserID
			usr.BlackPlayer = ""
		} else {
			usr.BlackPlayer = usr.UserID
			usr.WhitePlayer = ""
		}

		fmt.Println("\nGAME MODE SET\n", usr)
	}
}

func ClearOnlineGame(user *User) {
	var ok bool
	var opponent *User
	if user.WhitePlayer == user.UserID {
		opponent, ok = clients[user.BlackPlayer]
	} else {
		opponent, ok = clients[user.WhitePlayer]
	}
	if ok {
		user.WhitePlayer = ""
		user.BlackPlayer = ""
		opponent.WhitePlayer = ""
		opponent.BlackPlayer = ""
		fmt.Println("Cleared game opponent data")
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	usr, _ := GetUserAndGame(r)

	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("src/website/index.html")
		if err != nil {
			log.Print("Error parsing template: ", err)
		}
		err = t.Execute(w, usr)
		if err != nil {
			log.Print("Error during executing: ", err)
		}
		// if usr.WhitePlayer != "" {
		// 	ClearOnlineGame(usr)
		// }
	case "POST":
		r.ParseMultipartForm(0)
		message := r.FormValue("option")
		fmt.Println("Client in HomePage POST:", message)
		num, err := strconv.Atoi(string(message[0]))
		if err != nil {
			log.Print("Error with message from Client: ", err)
		}
		choseToBeWhite := true
		if string(message[1]) == "b" {
			choseToBeWhite = false
		}
		fmt.Println("CHOSE TO BE WHITE: ", choseToBeWhite, string(message[1]))
		if num == 2 {
			fmt.Println("Find Opponent")
			SetGameModeAndColour(usr.UserID, num, choseToBeWhite)
			session <- usr
			opp := <-usr.findOpponent
			fmt.Println("These are the users:", usr, opp)
		} else {
			SetGameModeAndColour(usr.UserID, num, choseToBeWhite)
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
// "atb" to accept takeback
// "rtb" to reject takeback
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
		}).ParseFiles("src/website/game.html")
		if err != nil {
			log.Print("Error parsing template: ", err)
		}
		fmt.Println(r.RemoteAddr, usr.UserID == usr.WhitePlayer)
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
			var lastMove string
			//check for en passant
			if result && game.Board.Board[destX][destY].Symbol == "P" && abs(startY-destY) == 1 && enPassantCheckPiece == " " {
				fmt.Fprintf(w, "Result:"+"enpassant"+strconv.Itoa(startX)+strconv.Itoa(destY)+checkText)
				lastMove = "Result:" + "enpassant" + strconv.Itoa(startX) + strconv.Itoa(destY) + checkText
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
				lastMove = "Result:" + "castle" + rookLocation + checkText
			} else if game.Board.promotePawn && usr.GameMode != 1 && usr.Game.IsWhiteTurn == true {
				game.Board.promotePawn = false
				fmt.Println("CHANGED PAWN PROMOTE BACK TO: ", game.Board.promotePawn)
				fmt.Fprintf(w, "Result:"+"pwn"+checkText)
				lastMove = "Result:" + "pwn" + checkText
			} else {
				fmt.Println("Result:" + strconv.FormatBool(result) + checkText)
				fmt.Fprintf(w, "Result:"+strconv.FormatBool(result)+checkText)
				lastMove = "Result:" + strconv.FormatBool(result) + checkText
			}
			if game.Board.promotePawn {
				game.Board.promotePawn = false
			}
			if result && usr.GameMode == 2 {
				fmt.Println("Sent")
				usr.opponenetMove <- lastMove
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
				usr.undoMove <- &Move{}
				// for opp.undoMoveRequested == "" {}
				move1 := <-usr.undoMove
				move2 := <-usr.undoMove

				// opp.undoMoveRequested = ""
				// if usr.undoMoveRequested == "rtb" {
				if *move1 == (Move{}) {
					fmt.Fprintf(w, "reject")
					// usr.undoMoveRequested = ""
					return
				}
				// usr.undoMoveRequested = ""
				fmt.Println(r.RemoteAddr, "take back accepted by opponent")
				// move := usr.lastMadeMove
				// move2 := usr.secondLastMove
				fmt.Fprintf(w, "true"+strconv.Itoa(move1.From.x)+strconv.Itoa(move1.From.y)+string(move1.From.Symbol)+string(strconv.FormatBool(move1.From.IsBlack)[0])+strconv.Itoa(move1.To.x)+strconv.Itoa(move1.To.y)+string(move1.To.Symbol)+string(strconv.FormatBool(move1.To.IsBlack)[0])+strconv.Itoa(move2.From.x)+strconv.Itoa(move2.From.y)+string(move2.From.Symbol)+string(strconv.FormatBool(move2.From.IsBlack)[0])+strconv.Itoa(move2.To.x)+strconv.Itoa(move2.To.y)+string(move2.To.Symbol)+string(strconv.FormatBool(move2.To.IsBlack)[0])+"ac")
				// usr.lastMadeMove = &Move{}
				// usr.secondLastMove = &Move{}
				return
			}
			result, move := game.undoTurn()
			if result {
				fmt.Println("Un did this move:", move)
				if usr.GameMode == 1 {
					result2, move2 := game.undoTurn()
					if result2 {
						fmt.Fprintf(w, "true"+strconv.Itoa(move.From.x)+strconv.Itoa(move.From.y)+string(move.From.Symbol)+string(strconv.FormatBool(move.From.IsBlack)[0])+strconv.Itoa(move.To.x)+strconv.Itoa(move.To.y)+string(move.To.Symbol)+string(strconv.FormatBool(move.To.IsBlack)[0])+strconv.Itoa(move2.From.x)+strconv.Itoa(move2.From.y)+string(move2.From.Symbol)+string(strconv.FormatBool(move2.From.IsBlack)[0])+strconv.Itoa(move2.To.x)+strconv.Itoa(move2.To.y)+string(move2.To.Symbol)+string(strconv.FormatBool(move2.To.IsBlack)[0])+"ai")
					} else {
						fmt.Fprintf(w, "false")
					}
				} else {
					fmt.Fprintf(w, "true"+strconv.Itoa(move.From.x)+strconv.Itoa(move.From.y)+string(move.From.Symbol)+string(strconv.FormatBool(move.From.IsBlack)[0])+strconv.Itoa(move.To.x)+strconv.Itoa(move.To.y)+string(move.To.Symbol)+string(strconv.FormatBool(move.To.IsBlack)[0]))
				}

			} else {
				fmt.Fprintf(w, "false")
			}
		case "col":
			fmt.Println("User:", usr)
			fmt.Println("Gamemode:", usr.GameMode)
			if usr.GameMode == 2 {
				fmt.Println("User:", usr)
				fmt.Println("response:", strconv.FormatBool(usr.WhitePlayer == usr.UserID))
				fmt.Fprintf(w, strconv.FormatBool(usr.WhitePlayer == usr.UserID)+strconv.Itoa(usr.GameMode))
			} else if usr.GameMode == 1 {
				//change later
				fmt.Fprintf(w, strconv.FormatBool(usr.WhitePlayer == usr.UserID)+strconv.Itoa(usr.GameMode))
			} else {
				fmt.Fprintf(w, "true"+strconv.Itoa(usr.GameMode))
			}
		case "opp":
			// for {
			// 	val, ok := usr.Game.Moves.Peek();
			// 	if ok && val.From.Symbol != " " && val.From.IsBlack == !game.IsWhiteTurn {
			// 		return
			// 	}
			// }
			fmt.Println("GOT OPP REQUEST FROM", r.RemoteAddr)
			fmt.Println("Current Turn:", game.IsWhiteTurn, usr.UserID == usr.WhitePlayer, usr.UserID == usr.BlackPlayer)
			// for game.IsWhiteTurn != (usr.UserID == usr.WhitePlayer) {}
			opp := GetOpponent(usr)
			if opp == nil {
				fmt.Fprintf(w, "false")
				return
			}
			var lastMove string
			select {
			case lastMove = <-opp.opponenetMove:
				fmt.Println("This is the move stackL", usr.Game.Moves)
				fmt.Println("This is the last move:", lastMove)
				fmt.Println("This is the user:", usr)
				fmt.Println("This is the opp:", opp)
				val, ok := usr.Game.Moves.Peek()
				fmt.Println("This is the val and ok:", val, ok)
				if ok {
					fmt.Fprintf(w, lastMove+strconv.Itoa(val.From.x)+strconv.Itoa(val.From.y)+strconv.Itoa(val.To.x)+strconv.Itoa(val.To.y))
				} else {
					fmt.Fprintf(w, "false")
				}
				fmt.Println("This is the last move:", lastMove+strconv.Itoa(val.From.x)+strconv.Itoa(val.From.y)+strconv.Itoa(val.To.x)+strconv.Itoa(val.To.y))
			case <-opp.undoMove:
				fmt.Println("Asked ", r.RemoteAddr, "for modal")
				fmt.Fprintf(w, "bck")
			}
			// for (opp.lastMove == "" || strings.Contains(opp.lastMove, "false")) && opp.undoMoveRequested == "" {}
		case "aim":
			fmt.Println("GOT AIM REQUEST FROM", r.RemoteAddr)
			start := time.Now()
			moveCounter := 0
			// for usr.Game.IsWhiteTurn == true {}
			for usr.Game.IsWhiteTurn == (usr.UserID == usr.WhitePlayer) {
			}
			AIMove := usr.Game.FindBestMove(!(usr.UserID == usr.WhitePlayer), 3, &moveCounter)
			usr.Game.makeMove(AIMove.From.x, AIMove.From.y, AIMove.To.x, AIMove.To.y)
			fmt.Println("Time Taken:", time.Since(start), "seconds")
			fmt.Fprintf(w, "Result:true"+getCheckMessage(usr.Game, true)+strconv.Itoa(AIMove.From.x)+strconv.Itoa(AIMove.From.y)+strconv.Itoa(AIMove.To.x)+strconv.Itoa(AIMove.To.y))
		case "atb":
			fmt.Println(r.RemoteAddr + "accepted")
			opp := GetOpponent(usr)
			//Assumes there are at least two moves on the move stack
			_, move := game.undoTurn()
			_, move2 := game.undoTurn()
			// opp.lastMadeMove = &move
			// opp.secondLastMove = &move2
			fmt.Println("These are the moves", move, move2)
			// opp.undoMoveRequested = "atb"
			// usr.undoMoveRequested = "answered"
			opp.undoMove <- &move
			opp.undoMove <- &move2

			fmt.Fprintf(w, "true"+strconv.Itoa(move.From.x)+strconv.Itoa(move.From.y)+string(move.From.Symbol)+string(strconv.FormatBool(move.From.IsBlack)[0])+strconv.Itoa(move.To.x)+strconv.Itoa(move.To.y)+string(move.To.Symbol)+string(strconv.FormatBool(move.To.IsBlack)[0])+strconv.Itoa(move2.From.x)+strconv.Itoa(move2.From.y)+string(move2.From.Symbol)+string(strconv.FormatBool(move2.From.IsBlack)[0])+strconv.Itoa(move2.To.x)+strconv.Itoa(move2.To.y)+string(move2.To.Symbol)+string(strconv.FormatBool(move2.To.IsBlack)[0])+"ac")
		case "rtb":
			// fmt.Println(r.RemoteAddr + "rejected")
			// opp := GetOpponent(usr)
			// opp.undoMoveRequested = "rtb"
			// usr.undoMoveRequested = "answered"
			opp := GetOpponent(usr)
			opp.undoMove <- &Move{}
			opp.undoMove <- &Move{}
			fmt.Fprintf(w, "reject")
		default:
			log.Print("HTTP Request Error")
		}

		// respond to client's request
		// fmt.Fprintf(w, "Server: %s \n", message+" | "+time.Now().Format(time.RFC3339))
	}
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

func GetPort() string {
	port := os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	go pairPlayers()
	// StartGame()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/website"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/game", GamePage)

	log.Fatal(http.ListenAndServe(GetPort(), nil))
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
