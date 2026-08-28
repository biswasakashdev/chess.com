package chess

type Turn int

const (
	TurnWhite Turn = Turn(Black)
	TurnBlack Turn = Turn(White)
)

type OutCome int

const (
	NoOutcome OutCome = iota
	BlackWin
	WhiteWin
)

type Chess interface {
	GetBoard() string
	MakeMove(string) error
	GetTurn() Turn
	GetOutCome() OutCome
}
