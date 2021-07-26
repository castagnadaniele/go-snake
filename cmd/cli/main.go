package main

import (
	"log"
	"time"

	"github.com/castagnadaniele/go-snake"
	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	err = screen.Init()
	if err != nil {
		log.Fatal(err)
	}
	width, height := screen.Size()

	s := snake.NewSnake(width, height)
	food := snake.NewFood(width, height)
	cloak := snake.NewCloak()
	defer cloak.Stop()
	game := snake.NewGame(s, cloak, food)
	view := snake.NewView(screen)
	defer view.Release()
	controller := snake.NewController(game, view)

	go controller.Start(time.Millisecond * 200)

	<-controller.WaitForQuitSignal()
}
