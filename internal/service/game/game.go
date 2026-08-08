package game

import "fmt"

type Position struct {
	Row int // 0 to 7
	Col int // 0 to 7
}

type Move struct {
	From Position
	To   Position
}

type Game struct {
	Board     Board
	Turn      Color
	FullMoves int
}

func NewGame() *Game {

	board := newBoard()
	return &Game{
		Board:     board,
		Turn:      White,
		FullMoves: 0,
	}
}

func PrintGame(g *Game) {

	fmt.Println("Board....")
	for _, arr := range g.Board {
		fmt.Println(arr)
	}

	fmt.Print("Turn...: ")
	if g.Turn == White {
		fmt.Print("While")
	} else {
		fmt.Print("Black")
	}

	fmt.Println()

	fmt.Println("Number of full moves", g.FullMoves)

}

func (g *Game) MakeMove(m Move) {

	// 1. Get the piece from the starting position
	piece := g.Board[m.From.Row][m.From.Col]

	// 2. Clear the starting position
	g.Board[m.From.Row][m.From.Col] = Piece{Type: None, Color: Nothing}

	// 3. Place the piece at the destination (this automatically overwrites/captures any piece there)
	g.Board[m.To.Row][m.To.Col] = piece

	// 4. Switch turns
	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
		g.FullMoves++
	}

}

func (g *Game) IsValidMove(m Move) bool {
	piece := g.Board[m.From.Row][m.From.Col]

	// Rule 1: You can't move an empty square
	if piece.Type == None {
		return false
	}

	// Rule 2: You can only move your own color's pieces
	if piece.Color != g.Turn {
		return false
	}

	// Rule 3: You can't capture your own piece
	destPiece := g.Board[m.To.Row][m.To.Col]
	if destPiece.Color == g.Turn {
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
