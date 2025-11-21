package domain

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         *User
}
