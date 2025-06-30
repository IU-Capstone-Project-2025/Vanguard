package models

type GetPlayersResponse struct {
	SessionCode string   `json:"sessionCode"`
	Players     []string `json:"players"`
}
