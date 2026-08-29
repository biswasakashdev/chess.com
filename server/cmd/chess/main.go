package main

import (
	"fmt"

	"github.com/biswasakashdev/chess.com/server/internal/chess"
)

func main() {

	game := chess.NewGame()

	// Get board works fine
	board := game.GetBoard()
	fmt.Println(board)

	cells := game.GetAvailableMoves("b8")
	fmt.Println(cells)
}
