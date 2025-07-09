package models

type ValidateCodeReq struct {
	Code     string `json:"code"`
	UserName string `json:"userName"`
}
