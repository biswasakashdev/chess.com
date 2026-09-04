package dtos

type GameDetails struct {
	GameId      string `json:"id"`
	WhitePlayer string `json:"white_player"`
	BlackPlayer string `json:"black_player"`
	Turn        string `json:"turn"`
	Board       string `json:"board"`
}
