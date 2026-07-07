package model

type UserInfoResponse struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
