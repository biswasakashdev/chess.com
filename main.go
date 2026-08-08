package main

import "github.com/biswasakashdev/chess.com/internal/service/game"

func main() {

	newGame := game.NewGame()

	game.PrintGame(newGame)

	newGame.MakeMove(game.Move{
		From: game.Position{
			Row: 6,
			Col: 2,
		},
		To: game.Position{
			Row: 4,
			Col: 2,
		},
	})

	game.PrintGame(newGame)
}
