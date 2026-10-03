package output

type CodeVerificationSender interface {
	SendCodeVerification(to string, code string) error
}