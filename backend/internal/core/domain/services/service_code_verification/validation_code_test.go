package service_code_verification_test

import (
	"testing"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_code_verification"
)

func TestCodeValidator(t *testing.T) {
	validator := service_code_verification.CodeValidator{}
	cases := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"wrong type", 123456, true},
		{"empty", "", true},
		{"too short", "12345", true},
		{"too long", "1234567", true},
		{"negative", "-12345", true},
		{"contains letters", "12a456", true},
		{"valid", "012345", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate(%v) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
		})
	}
}
