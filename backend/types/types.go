package types

type Map map[string]any

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserSignupDTO struct {
	Username string `validate:"required,min=3,max=20,alphanum"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=50"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt int    `json:"created_at"`
	AvatarURL string `json:"avatar_url"`
}

type UserSession struct {
	ID                   string
	UserID               string
	CreatedAt            int
	ExpiresAt            int
	RefreshToken         string
	AccessToken          string
	AccessTokenExpiresAt string
}

type UserLoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=50"`
}
