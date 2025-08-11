package usecase

type AuthPrincipal struct {
	ID    uint
	Role  string
	Email string
}

type HashManager interface {
	HashPassword(password string) (string, error)
	Compare(password string, hash string) error
}

type TokenManager interface {
	Generate(userID uint, role, email string) (string, error)
	Verify(tokenString string) (*AuthPrincipal, error)
}
