package game

type MsgLogin struct {
	MsgId    int    `json:"msgid"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   int64  `json:"userid"`
}
