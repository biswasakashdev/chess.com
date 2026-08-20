package chess

type PieceType int

const Size int = 8

const (
	None PieceType = iota // Nothing on the cell
	Rook
	Knight
	Bishop
	King
	Queen
	Pawn
)

type Color int

func (cl Color) String() string {
	if cl == Black {
		return "BLACK"
	}

	return "WHITE"
}

const (
	Nothing Color = iota
	Black
	White
)

type Piece struct {
	Type  PieceType
	Color Color
}

func (p Piece) String() string {

	// When the color black then this execute
	if p.Color == Black {
		return []string{".", "r", "n", "b", "k", "q", "p"}[p.Type]
	}

	return []string{".", "R", "N", "B", "K", "Q", "P"}[p.Type]
}

type Board [8][8]Piece

func newBoard() Board {

	var grid Board

	defaultPiece := Piece{
		Color: Nothing,
		Type:  None,
	}

	for row := range grid {
		for col := range grid[row] {
			grid[row][col] = defaultPiece
		}
	}

	// Set rooks
	grid[0][0] = Piece{Color: Black, Type: Rook}
	grid[0][7] = Piece{Color: Black, Type: Rook}
	grid[7][0] = Piece{Color: White, Type: Rook}
	grid[7][7] = Piece{Color: White, Type: Rook}

	//Set knights

	grid[0][1] = Piece{Color: Black, Type: Knight}
	grid[0][6] = Piece{Color: Black, Type: Knight}
	grid[7][1] = Piece{Color: White, Type: Knight}
	grid[7][6] = Piece{Color: White, Type: Knight}

	// Set bishops

	grid[0][2] = Piece{Color: Black, Type: Bishop}
	grid[0][5] = Piece{Color: Black, Type: Bishop}
	grid[7][2] = Piece{Color: White, Type: Bishop}
	grid[7][5] = Piece{Color: White, Type: Bishop}

	// Set queen

	grid[0][3] = Piece{Color: Black, Type: Queen}
	grid[7][3] = Piece{Color: White, Type: Queen}

	// Set King

	grid[0][4] = Piece{Color: Black, Type: King}
	grid[7][4] = Piece{Color: White, Type: King}
	// Set
	for col := range Size {
		grid[1][col] = Piece{
			Color: Black,
			Type:  Pawn,
		}
	}

	for col := range Size {
		grid[6][col] = Piece{
			Type:  Pawn,
			Color: White,
		}
	}

	return grid

}

// Create a res

func (b Board) ToBoardRes() [8][8]string {

	var board [8][8]string
	for row := range b {
		for col := range board[row] {
			board[row][col] = b[row][col].String()
		}
	}

	return board
}
