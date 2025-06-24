package models

type SessionCreateResponse struct {
	Jwt              string `json:"jwt"`
	ServerWsEndpoint string `json:"serverWsEndpoint"`
	SessionId        string `json:"sessionId"`
}
