package dtos

type GameTurn string

var (
	GameTurnWhite GameTurn = "white"
	GameTurnBlack GameTurn = "black"
)

type GameDetails struct {
	GameId      string   `json:"id"`
	WhitePlayer string   `json:"white_player"`
	BlackPlayer string   `json:"black_player"`
	Turn        GameTurn `json:"turn"`
	Board       string   `json:"board"`
}
