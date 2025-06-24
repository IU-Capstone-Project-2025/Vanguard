package models

type ValidateCodeReq struct {
	Code   string `json:"code"`
	UserId string `json:"userId"`
}
