package chess

type Turn int

const (
	TurnWhite Turn = Turn(Black)
	TurnBlack Turn = Turn(White)
)

type OutCome int

const (
	NoOutCome OutCome = iota
	BlackWin
	WhiteWin
)

type Chess interface {
	GetBoard() string
	MakeMove(string) error
	GetAvailableMoves(string) string
	GetTurn() Turn
	GetOutCome() OutCome
}
