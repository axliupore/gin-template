package request

// UserRegister 用户注册结
type UserRegister struct {
	Account       string `json:"account"`
	Password      string `json:"password"`
	CheckPassword string `json:"checkPassword"`
}

// UserLogin 用户登录
type UserLogin struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// UserUpdate 更新用户信息
type UserUpdate struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Profile  string `json:"profile"`
	Gender   int8   `json:"gender"`
}

// UserSearch 查找用户
type UserSearch struct {
	PageRequest
	Text string
}

type UserDelete struct {
	Id int64 `json:"id"`
}
