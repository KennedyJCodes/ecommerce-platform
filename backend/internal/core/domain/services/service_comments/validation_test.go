package service_comments_test

import (
	"testing"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_comments"
)

func TestCommentValidator(t *testing.T) {
	validator := service_comments.CommentValidator{}
	cases := []struct{
		name string
		input interface{}
		wantErr bool
	}{
		{"Wrong type", "is not CommentValidationData", true},
		{"Empty content", service_comments.CommentValidationData{Content: "", Rating: 5}, true},
		{"Rating not provided (0)", service_comments.CommentValidationData{Content: "I liked the watch", Rating: 0}, true},
		{"Rating out of range (less than 1)", service_comments.CommentValidationData{Content: "I liked the watch", Rating: -1}, true},
		{"Rating out of range (greater than 5)", service_comments.CommentValidationData{Content: "I liked the watch", Rating: 6}, true},
		{"Valid data", service_comments.CommentValidationData{Content: "I liked the watch", Rating: 5}, false},
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