package model

const (
	DBError               = "database_error"
	JwtSigningError       = "jwt_signing_error"
	FailedToParseJwtToken = "failed_to_parse_jwt_token"
	InvalidJwtToken       = "invalid_jwt_token"
	ExpiredJwtToken       = "expired_jwt_token"
	EmptyAccessToken      = "access_token_not_found"
	EmptyRefreshToken     = "access_token_not_found"
	InvalidPassword       = "invalid_password"
	TimeParseError        = "time_parse_error"
	BcryptHashError       = "bcrypt_hash_error"
)
