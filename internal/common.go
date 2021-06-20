package internal

import (
	"github.com/chiahsoon/go_scaffold/internal/auth"
	"github.com/chiahsoon/go_scaffold/internal/dals"
)

// DALs == Services

var AuthService auth.AuthService
var UserService dals.UserDAL

func Init() {
	AuthService = auth.AuthService{}
	UserService = dals.UserDAL{}
}
