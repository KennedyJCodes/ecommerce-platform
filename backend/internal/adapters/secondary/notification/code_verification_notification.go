package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type codeVerificationNotification struct {
	apiKey     string
	fromEmail  string
	httpClient *http.Client
}

func NewCodeVerificationNotification(apiKey, fromEmail string) *codeVerificationNotification {
	return &codeVerificationNotification{
		apiKey:     apiKey,
		fromEmail:  fromEmail,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type CodeVerificationPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func (m *codeVerificationNotification) SendCodeVerification(to string, code string) error {
	payload := CodeVerificationPayload{
		From:    m.fromEmail,
		To:      []string{to},
		Subject: "Your verification code",
		HTML:    fmt.Sprintf("<p>Your verification code is: <strong>%s</strong></p><p>It expires in 10 minutes.</p>", code),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializing payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
    	return fmt.Errorf("resend responded with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}