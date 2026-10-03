package service_code_verification_test

import (
	"testing"
	"unicode"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_code_verification"
)

func TestGenerateCodeVerification(t *testing.T) {
	service := service_code_verification.NewCodeVerificationService()
	validator := service_code_verification.CodeValidator{}

	for i := 0; i < 100; i++ {
		code, err := service.GenerateCodeVerification()
		if err != nil {
			t.Fatalf("GenerateCodeVerification() error = %v", err)
		}
		if len(code) != 6 {
			t.Fatalf("generated code length = %d, want 6", len(code))
		}
		for _, char := range code {
			if !unicode.IsDigit(char) {
				t.Fatalf("generated code contains non-digit character %q", char)
			}
		}
		if err := validator.Validate(code); err != nil {
			t.Fatalf("generated code failed validation: %v", err)
		}
	}
}
