package input

type CSRFService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenValue string, userID string) error
}
