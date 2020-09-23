package main

import (
	"fmt"
	// "io/ioutil"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	// "time"
)

var (
	game *Game
	username string
	users = make(map[string]*Game)
)

//get rid of this
type WebpageData struct {
	Game   *Game
	Player string
}

func GetGame(r *http.Request) {
	ip := r.RemoteAddr
	if string(ip[0]) == "[" {
		//parse localHost IP
		ip = ip[:5]//since localhost is like [::1]:PORT
	} else {
		ip = strings.Split(ip, ":")[0]//Splits IPv4 address and port into just address
	}
	if g, ok := users[ip]; ok {
		game = g
	} else {
		users[ip] = &Game{Board{}, true}
		game = users[ip]
		SetupBoard(&game.Board)
	}
	username = ip
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	GetGame(r)
	if r.Method != "GET" {
		log.Print("Not GET request recieved on Home Page")
		return
	}
	vars := WebpageData{game, username}
	t, err := template.ParseFiles("website/index.html")
	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, vars)
	if err != nil {
		log.Print("Error during executing: ", err)
	}
}

// Message will be either:
// "opt 00" for give move options for piece at (0,0)
// "mov 00 01" for move piece at (0,0) to (0,1)
// "pwn 00q" to promote pawn at 00 to a queen
// "rst " to restart the game
func GamePage(w http.ResponseWriter, r *http.Request) {
	GetGame(r)
	fmt.Println("IP:", r.RemoteAddr)
	fmt.Println(users)
	switch r.Method {
	case "GET":
		// game = Game{Board{}, true}
		// SetupBoard(&game.Board)
		tmpl, err := template.New("game.html").Funcs(template.FuncMap{
			"minus": func(a, b int) int {
				return a - b
			},
			"add": func(a, b int) int {
				return a + b
			},
		}).ParseFiles("website/game.html")
		// t, err := template.ParseFiles("website/game.html")
		if err != nil {
			log.Print("Error parsing template: ", err)
		}
		err = tmpl.Execute(w, WebpageData{game, username})
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

			checkText := getCheckMessage(result) 
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
				fmt.Fprintf(w, "Result:" + "castle" + rookLocation + checkText)
			} else if game.Board.promotePawn {
				game.Board.promotePawn = false
				fmt.Println("CHANGED PAWN PROMOTE BACK TO: ", game.Board.promotePawn)
				fmt.Fprintf(w, "Result:" + "pwn" + checkText)

			} else {
				fmt.Println("Result:" + strconv.FormatBool(result) + checkText)
				fmt.Fprintf(w, "Result:" + strconv.FormatBool(result) + checkText)
			}

		case "pwn":
			x, y := int(rest[0])-int('0'), int(rest[1])-int('0')
			fmt.Println("This is the string of the oanw to promote:", string(rest[2]))
			result := game.Board.Board[x][y].promotePawn(&game.Board, string(rest[2]))
			checkText := getCheckMessage(result) 
			fmt.Fprintf(w, strconv.FormatBool(result) + checkText)
		
		case "rst":
			SetupBoard(&game.Board)
			game.IsWhiteTurn = true
			fmt.Fprintf(w, "reload")
		case "ply":
			fmt.Fprintf(w, strconv.FormatBool(game.IsWhiteTurn) + getCheckMessage(true))
		default:
			log.Print("HTTP Request Error")
		}

		// respond to client's request
		// fmt.Fprintf(w, "Server: %s \n", message+" | "+time.Now().Format(time.RFC3339))
	}
}

func getCheckMessage(result bool) string{
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
