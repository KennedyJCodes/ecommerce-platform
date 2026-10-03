package dto

// VerifyCodeRequest contains the verification code submitted by the user.
type VerifyCodeRequest struct {
	Code string `json:"code"`
}
