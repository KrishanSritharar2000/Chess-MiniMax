package main

import (
	"fmt"
	// "io/ioutil"
	"net/http"
	"html/template"
	"log"
)

var game = Game{Board{}, true}

type WebpageData struct {
	Game   *Game
	Player string
}

func HomePage(w http.ResponseWriter, r *http.Request) {
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

func GamePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("website/game.html")
	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, WebpageData{&game, "Start Game"})
	if err != nil {
		log.Print("Error during executing: ", err)
	}
}

func GamePageSelected(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SUCCESS WE MADE IT TO SELECTED")
	t, err := template.ParseFiles("website/game.html")
	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, WebpageData{&game, "Returned form Selected"})
	if err != nil {
		log.Print("Error during executing: ", err)
	}}


func main() {
	// StartGame()
	SetupBoard(&game.Board)

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
