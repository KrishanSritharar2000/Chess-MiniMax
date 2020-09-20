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
	game = Game{Board{}, true}
)

type WebpageData struct {
	Game   *Game
	Player string
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Print("Not GET request recieved on Home Page")
		return
	}
	vars := WebpageData{&game, "Krishan"}
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
func GamePage(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		SetupBoard(&game.Board)
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
		err = tmpl.Execute(w, WebpageData{&game, "Start Game"})
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
			fmt.Println("Move the piece at:", startX, startY, game.Board.Board[startX][startY], "to:", destX, destY, game.Board.Board[destX][destY])
			result := game.makeMove(startX, startY, destX, destY)
			fmt.Fprintf(w, "Result:" + strconv.FormatBool(result))
		default:
			log.Print("HTTP Request Error")
		}

		// respond to client's request
		// fmt.Fprintf(w, "Server: %s \n", message+" | "+time.Now().Format(time.RFC3339))
	}
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

func GamePageSelected(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SUCCESS WE MADE IT TO SELECTED")
	r.ParseForm()
	fmt.Println("Form: ")
	fmt.Println(r.Form)
	t, err := template.ParseFiles("website/game.html")
	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, WebpageData{&game, "Returned form Selected"})
	if err != nil {
		log.Print("Error during executing: ", err)
	}
}

func main() {
	// StartGame()

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./website"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/game", GamePage)
	http.HandleFunc("/gameSelected", GamePageSelected)

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
