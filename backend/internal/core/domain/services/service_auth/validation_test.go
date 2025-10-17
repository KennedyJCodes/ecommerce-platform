package service_auth_test

import (
	"testing"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_auth"
)

func TestUserNameValidator(t *testing.T) {
	validator := &service_auth.UserNameValidator{}
	cases := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"Empty", "", true},
		{"TooShort", "abc", true},
		{"Exact", "abcde", false},
		{"Valid", "valid_username", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.input)
			if tc.wantErr && err == nil {
				t.Errorf("I wanted an error with input=%q, but there wasn't one.", tc.input)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("I didn't want an error with input=%q, but there was: %v", tc.input, err)
			}
		})
	}
}

func TestPasswordValidator(t *testing.T) {
	validator := &service_auth.PasswordValidator{}
	cases := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"Empty", "", true},
		{"TooShort", "A1!", true},
		{"NoUpper", "abcdefgh1@", true},
		{"NoDigit", "ABCDEFG!@---not", true},
		{"NoSpecial", "Abcdefghi1", true},
		{"Valid1", "abcdEFH#123---02not", false},
		{"Valid2", "THIS_password_is_validx2#", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.input)
			if tc.wantErr && err == nil {
				t.Errorf("I wanted an error with input=%q, but there wasn't one", tc.input)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("I didn't want an error with input=%q, but there was: %v", tc.input, err)
			}
		})
	}
}