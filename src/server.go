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
	}
}

type User struct {
	Game     *Game
	GameMode int //0 local PvP, 1 AI, 2 Online PvP
	userID string// ip address
	whitePlayer, blackPlayer string //ip address of players 
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
	clients[ip] = &User{&Game{Board{}, true, &MoveStack{}}, 0, ip, "", ""}
	game := clients[ip].Game
	SetupBoard(&game.Board)
	return clients[ip], game
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
		if usr.whitePlayer != "" {
			ClearOnlineGame(usr)
		}
	case "POST":
		r.ParseMultipartForm(0)
		message := r.FormValue("option")
		fmt.Println("Client:", message)
		num, err := strconv.Atoi(message)
		if err != nil {
			log.Print("Error with message from Client: ", err)
		}
		if num == 2 {
			fmt.Println("Find Opponent")
			SetGameMode(usr.whitePlayer, num)
			session <- usr
			for usr.whitePlayer == "" {}
		} else {
			SetGameMode(usr.whitePlayer, num)
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
		fmt.Println("Message from Client:", message)

		splitIndex := strings.Index(message, " ")
		if splitIndex == -1 {
			log.Print("Error in HTTP request:", message, len(message))
		}
		fmt.Println("message:", message, " splitIndex:", splitIndex)
		rest := message[(splitIndex + 1):]
		opcode := message[:splitIndex]

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
			} else if game.Board.promotePawn {
				game.Board.promotePawn = false
				fmt.Println("CHANGED PAWN PROMOTE BACK TO: ", game.Board.promotePawn)
				fmt.Fprintf(w, "Result:"+"pwn"+checkText)

			} else {
				fmt.Println("Result:" + strconv.FormatBool(result) + checkText)
				fmt.Fprintf(w, "Result:"+strconv.FormatBool(result)+checkText)
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
			fmt.Fprintf(w, strconv.FormatBool(game.IsWhiteTurn)+getCheckMessage(game, true))
		case "bck":
			time.Sleep(time.Second*5)
			result, move := game.undoTurn()
			if result {
				fmt.Println("Un did this move:", move)
				fmt.Fprintf(w, "true" + strconv.Itoa(move.From.x) + strconv.Itoa(move.From.y) + string(move.From.Symbol) + string(strconv.FormatBool(move.From.IsBlack)[0]) + strconv.Itoa(move.To.x) + strconv.Itoa(move.To.y) + string(move.To.Symbol) + string(strconv.FormatBool(move.To.IsBlack)[0]))
			} else {
				fmt.Fprintf(w, "false")
			}

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
