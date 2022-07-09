package dto

type ResLogin struct {
	Res
	Data Token `json:"data"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
	TokenType   string `json:"token_type"`
}
