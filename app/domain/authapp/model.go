package authapp

import "encoding/json"

type TokenResponse struct {
	Token string `json:"token"`
}

func (token TokenResponse) Encode() (data []byte, contentType string, err error) {
	data, err = json.Marshal(token)
	return data, "application/json", err
}
