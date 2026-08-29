package chess

import (
	"errors"
	"fmt"
	"strings"
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

// GetBoard returns the board state as a string where each row is separated by '\n'
func (g *Game) GetBoard() string {
	var rows []string
	for r := range 8 {
		var row string
		for _, val := range g.board[r] {
			row = row + " " + val.String()
		}
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n")
}

// GetTurn returns which color's turn it is
func (g *Game) GetTurn() Turn {
	return Turn(g.turn)
}

// GetOutCome checks if any king was captured and returns game outcome
func (g *Game) GetOutCome() OutCome {
	whiteKingFound := false
	blackKingFound := false

	for r := range 8 {
		for c := range 8 {
			if g.board[r][c].String() == "K" {
				whiteKingFound = true
			} else if g.board[r][c].String() == "k" {
				blackKingFound = true
			}
		}
	}

	if !whiteKingFound {
		return BlackWin
	}
	if !blackKingFound {
		return WhiteWin
	}
	return NoOutCome
}

// GetAvailableMoves returns a space-separated string of reachable targets from a cell (e.g., "e5 e6")
func (g *Game) GetAvailableMoves(cell string) string {
	r, c, err := parseCell(cell)
	if err != nil {
		return ""
	}

	moves := g.generateMoves(r, c)
	var targets []string
	for _, m := range moves {
		targets = append(targets, toCell(m[0], m[1]))
	}

	return strings.Join(targets, " ")
}

// MakeMove accepts moves in UCI format like "e2e4" or "e4a6"
func (g *Game) MakeMove(moveStr string) error {
	if len(moveStr) != 4 {
		return errors.New("invalid move format, expected 4 characters like 'e2e4'")
	}

	if g.GetOutCome() != NoOutCome {
		return errors.New("game is already over")
	}

	fromCell := moveStr[0:2]
	toCellStr := moveStr[2:4]

	fromR, fromC, err := parseCell(fromCell)
	if err != nil {
		return err
	}
	toR, toC, err := parseCell(toCellStr)
	if err != nil {
		return err
	}

	piece := g.board[fromR][fromC]
	if piece.isNothing() {
		return errors.New("no piece at source square")
	}

	// Turn check
	if (g.turn == White && !piece.isWhite()) || (g.turn == Black && !piece.isBlack()) {
		return errors.New("not your piece to move")
	}

	// Validate against piece's legal pseudo-moves
	available := g.generateMoves(fromR, fromC)
	valid := false
	for _, m := range available {
		if m[0] == toR && m[1] == toC {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("illegal move for this piece")
	}

	// Execute move
	g.board[toR][toC] = piece
	g.board[fromR][fromC] = Piece{
		Color: Nothing,
		Type:  None,
	}

	// Auto-promote pawns to Queen on reaching last rank
	if piece.String() == "P" && toR == 0 {
		g.board[toR][toC] = Piece{
			Type:  Queen,
			Color: White,
		}
	} else if piece.String() == "p" && toR == 7 {
		g.board[toR][toC] = Piece{
			Type:  Queen,
			Color: Black,
		}
	}

	// Switch turn and update move counters
	if g.turn == White {
		g.turn = Black
	} else {
		g.turn = White
		g.fullMoves++
	}

	return nil
}

// --- Movement Logic & Path Generation ---

func (g *Game) generateMoves(r, c int) [][2]int {
	piece := g.board[r][c]
	if piece.isNothing() {
		return nil
	}

	var moves [][2]int
	white := piece.isWhite()
	lower := strings.ToLower(piece.String())

	switch lower {
	case "p":
		dir := -1 // White moves up (rank 2 to rank 8 -> index 6 down to 0)
		startRow := 6
		if !white {
			dir = 1 // Black moves down
			startRow = 1
		}

		// 1 step forward
		nextR := r + dir
		if inBounds(nextR, c) && g.board[nextR][c].isNothing() {
			moves = append(moves, [2]int{nextR, c})
			// 2 steps forward from starting row
			if r == startRow && g.board[r+2*dir][c].isNothing() {
				moves = append(moves, [2]int{r + 2*dir, c})
			}
		}

		// Diagonal captures
		for _, dc := range []int{-1, 1} {
			capC := c + dc
			if inBounds(nextR, capC) {
				target := g.board[nextR][capC]
				if !target.isNothing() && piece.isEnemy(target) {
					moves = append(moves, [2]int{nextR, capC})
				}
			}
		}

	case "n":
		offsets := [][2]int{
			{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
			{1, -2}, {1, 2}, {2, -1}, {2, 1},
		}
		for _, off := range offsets {
			nr, nc := r+off[0], c+off[1]
			if inBounds(nr, nc) && (g.board[nr][nc].isNothing() || piece.isEnemy(g.board[nr][nc])) {
				moves = append(moves, [2]int{nr, nc})
			}
		}

	case "b":
		moves = append(moves, g.raycast(r, c, piece, [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}})...)

	case "r":
		moves = append(moves, g.raycast(r, c, piece, [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}})...)

	case "q":
		dirs := [][2]int{
			{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
			{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		}
		moves = append(moves, g.raycast(r, c, piece, dirs)...)

	case "k":
		dirs := [][2]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1}, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		}
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if inBounds(nr, nc) && (g.board[nr][nc].isNothing() || piece.isEnemy(g.board[nr][nc])) {
				moves = append(moves, [2]int{nr, nc})
			}
		}
	}

	return moves
}

func (g *Game) raycast(r, c int, piece Piece, directions [][2]int) [][2]int {
	var moves [][2]int
	for _, d := range directions {
		currR, currC := r+d[0], c+d[1]
		for inBounds(currR, currC) {
			target := g.board[currR][currC]
			if target.isNothing() {
				moves = append(moves, [2]int{currR, currC})
			} else {
				if piece.isEnemy(target) {
					moves = append(moves, [2]int{currR, currC})
				}
				break // ray blocked
			}
			currR += d[0]
			currC += d[1]
		}
	}
	return moves
}

// --- Coordinate Helpers ---

func parseCell(cell string) (int, int, error) {
	if len(cell) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinate %s", cell)
	}
	file := cell[0]
	rank := cell[1]

	if file < 'a' || file > 'h' || rank < '1' || rank > '8' {
		return 0, 0, fmt.Errorf("out of range coordinate %s", cell)
	}

	c := int(file - 'a')
	r := 8 - int(rank-'0')
	return r, c, nil
}

func toCell(r, c int) string {
	file := string(rune('a' + c))
	rank := string(rune('0' + (8 - r)))
	return file + rank
}

func inBounds(r, c int) bool {
	return r >= 0 && r < 8 && c >= 0 && c < 8
}
