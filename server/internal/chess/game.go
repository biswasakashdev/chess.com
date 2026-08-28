package chess

import (
	"fmt"
)

// JSON Payload structures
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Move struct {
	From Position
	To   Position
}

type BoardResponse struct {
	Board [8][8]string `json:"board"`
	Turn  string       `json:"turn"`
}

// Game session

type Game struct {
	board     Board
	turn      Color
	fullMoves int
}

func NewGame() *Game {
	board := newBoard()
	return &Game{
		board:     board,
		turn:      White,
		fullMoves: 0,
	}
}

func (g *Game) MakeMove(moveUCI string) error {

	move := getMove(moveUCI)

	// 1. Get the piece from the starting position
	piece := g.board[move.From.Row][move.From.Col]

	// 2. Clear the starting position
	g.board[move.From.Row][move.From.Col] = Piece{Type: None, Color: Nothing}

	// 3. Place the piece at the destination (this automatically overwrites/captures any piece there)
	g.board[move.To.Row][move.To.Col] = piece

	// 4. Switch turns
	if g.turn == White {
		g.turn = Black
	} else {
		g.turn = White
		g.fullMoves++
	}

	return nil

}

func getMove(_moveUCI string) Move {
	return Move{
		From: Position{
			Row: 1,
		},
	}
}

// Return the State of the board
func (g *Game) GetBoard() string {

	var board string
	for row := range g.board {
		for col := range board[row] {
			board = fmt.Sprintf("%s%s", board, g.board[row][col])
		}
		board = fmt.Sprintf("%s\n", board)
	}
	return board
}

func (g *Game) isMoveValid(m Move) bool {
	piece := g.board[m.From.Row][m.From.Col]

	// Rule 1: You can't move an empty square
	if piece.Type == None {
		return false
	}

	// Rule 2: You can only move your own color's pieces
	if piece.Color != g.turn {
		return false
	}

	// Rule 3: You can't capture your own piece
	destPiece := g.board[m.To.Row][m.To.Col]
	if destPiece.Color == g.turn {
		return false
	}

	// Rule 4: Piece-specific logic
	switch piece.Type {
	case Pawn:
		return validatePawnMove(m)
	case Knight:
		return validateKnightMove(m)
		// Add other pieces...
	}

	return true
}

func validatePawnMove(_ Move) bool {
	return true
}

func validateKnightMove(_ Move) bool {
	return true
}

func (g *Game) GetTurn() Turn {

	return TurnBlack
}

func (g *Game) GetOutCome() OutCome {

	return NoOutcome
}
