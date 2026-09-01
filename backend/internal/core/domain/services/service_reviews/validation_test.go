package service_reviews_test

import (
	"testing"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_reviews"
)

func TestReviewValidator(t *testing.T) {
	validator := service_reviews.ReviewValidator{}
	cases := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"Wrong type", "is not ReviewValidationData", true},
		{"Empty content", service_reviews.ReviewValidationData{Content: "", Rating: 5}, true},
		{"Rating not provided (0)", service_reviews.ReviewValidationData{Content: "I liked the watch", Rating: 0}, true},
		{"Rating out of range (less than 1)", service_reviews.ReviewValidationData{Content: "I liked the watch", Rating: -1}, true},
		{"Rating out of range (greater than 5)", service_reviews.ReviewValidationData{Content: "I liked the watch", Rating: 6}, true},
		{"Valid data", service_reviews.ReviewValidationData{Content: "I liked the watch", Rating: 5}, false},
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
