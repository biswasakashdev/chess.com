package dtos

type GameTurn string

var (
	GameTurnWhite GameTurn = "white"
	GameTurnBlack GameTurn = "black"
)

type GameMove struct {
	UserId string `json:"id"`
	Move   string `json:"move"`
}

type GameDetails struct {
	GameId      string     `json:"id"`
	WhitePlayer UserResp   `json:"white_player"`
	BlackPlayer UserResp   `json:"black_player"`
	Turn        GameTurn   `json:"turn"`
	Board       string     `json:"board"`
	History     []GameMove `json:"history"`
}

type AvailableMovesResp struct {
	Moves string `json:"moves"`
}
