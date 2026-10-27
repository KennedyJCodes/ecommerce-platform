package output

import "context"

type CodeVerificationSender interface {
	SendCodeVerification(ctx context.Context, to string, code string) error
}
