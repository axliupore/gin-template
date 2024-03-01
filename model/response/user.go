package response

import "github.com/axliupore/gin-template/model"

type UserLoginResponse struct {
	User      model.User `json:"user"`
	Token     string     `json:"token"`
	ExpiresAt int64      `json:"expiresAt"`
}
