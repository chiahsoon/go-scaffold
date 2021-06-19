package models

type UserRefreshToken struct {
	User         User // Do not assign the value of User when doing DB operations, use UserID
	UserID       string
	RefreshToken string `json:"refresh_token"`
	CreatedAt    uint   `json:"created_at"`
	UpdatedAt    uint   `json:"updated_at"`
	DeletedAt    *uint  `json:"deleted_at"`
}
