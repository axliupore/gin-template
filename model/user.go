package model

import (
	"time"
)

type User struct {
	Id       int64  `db:"id" json:"id"`
	Account  string `db:"account" json:"account"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
	Avatar   string `db:"avatar" json:"avatar"`
	Email    string `db:"email" json:"email"`
	Phone    string `db:"phone" json:"phone"`
	Profile  string `db:"profile" json:"profile"`
	Gender   int8   `db:"gender" json:"gender"`
	Role     string `db:"role" json:"role"`
	Model
}

func NewUser(account, password string) *User {
	user := &User{
		Account:  account,
		Password: password,
	}
	user.init()
	return user
}

// Initialize 用于初始化用户模型的默认值
func (u *User) init() {
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()
	u.IsDelete = 0
	u.Role = "user"
}
