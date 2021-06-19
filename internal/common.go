package internal

import (
	"github.com/chiahsoon/go_scaffold/internal/auth"
	"github.com/chiahsoon/go_scaffold/internal/dals"
)

// DALs == Services

var AuthService auth.AuthService
var UserService dals.UserDAL
var UserRefreshTokenService dals.UserRefreshTokenDAL

func Init() {
	AuthService = auth.AuthService{}
	UserService = dals.UserDAL{}
	UserRefreshTokenService = dals.UserRefreshTokenDAL{}
}
