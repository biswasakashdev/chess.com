package game

// JSON Payload structures
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Move struct {
	From Position `json:"from"`
	To   Position `json:"to"`
}

type BoardResponse struct {
	Board [8][8]string `json:"board"`
	Turn  string       `json:"turn"`
}

// Geme session

type Game struct {
	Board     Board
	Turn      Color
	FullMoves int
}

func newGame() *Game {

	board := newBoard()
	return &Game{
		Board:     board,
		Turn:      White,
		FullMoves: 0,
	}
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
