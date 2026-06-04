package main

import "fmt"

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

const (
	Nothing Color = iota
	Black
	White
)

type Piece struct {
	Type  PieceType
	Color Color
}

func (pt PieceType) String() string {
	return []string{".", "R", "N", "B", "K", "Q", "P"}[pt]
}

func (p Piece) String() string {

	// When the color black then this execute
	if p.Color == Black {
		return []string{".", "r", "n", "b", "k", "q", "p"}[p.Type]
	}
	// Other wise call the above function.
	return p.Type.String()
}

type Board [][]Piece

func NewBoard() Board {

	var grid [][]Piece = make([][]Piece, Size)

	for i := range grid {
		curr := make([]Piece, Size)
		// Initialize all the blocks with Empty Piece(No color and No piece).
		for j := range curr {
			curr[j] = Piece{
				Type:  None,
				Color: Nothing,
			}
		}
		grid[i] = curr
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

	// Set up

	grid[0][4] = Piece{Color: Black, Type: King}
	grid[7][4] = Piece{Color: White, Type: King}

	for col := range Size {
		grid[1][col] = Piece{
			Color: White,
			Type:  Pawn,
		}
	}

	for col := range Size {
		grid[6][col] = Piece{
			Type:  Pawn,
			Color: Black,
		}
	}

	return grid

}

func main() {

	board := NewBoard()

	for col := range board {
		fmt.Println(board[col])
	}

}
