package input

type CodeVerificationService interface {
	GenerateCodeVerification() (string, error)
}