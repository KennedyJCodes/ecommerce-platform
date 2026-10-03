package service_code_verification

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type CodeVerificationService struct{}

func NewCodeVerificationService() *CodeVerificationService {
	return &CodeVerificationService{}
}

func (s *CodeVerificationService) GenerateCodeVerification() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("error generate code: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}