package cookies_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

func TestWithValue(t *testing.T) {
	config := cookies.CookieConfig{}

	option := cookies.WithValue("test-value")
	option(&config)

	if config.Value != "test-value" {
		t.Errorf("WithValue did not set the value correctly. Expected: %s, Got: %s", "test-value", config.Value)
	}
}

func TestWithMaxAge(t *testing.T) {
	config := cookies.CookieConfig{}

	expectedDuration := 1 * time.Hour
	option := cookies.WithMaxAge(expectedDuration)
	option(&config)

	if config.MaxAge != expectedDuration {
		t.Errorf("WithMaxAge did not set the duration correctly. Expected: %v, Got: %v", expectedDuration, config.MaxAge)
	}
}

func TestWithHttpOnly(t *testing.T) {
	config := cookies.CookieConfig{}

	option := cookies.WithHttpOnly(true)
	option(&config)

	if !config.HttpOnly {
		t.Errorf("WithHttpOnly did not set the flag correctly. Expected: %v, Got: %v", true, config.HttpOnly)
	}

	option = cookies.WithHttpOnly(false)
	option(&config)

	if config.HttpOnly {
		t.Errorf("WithHttpOnly did not set the flag correctly. Expected: %v, Got: %v", false, config.HttpOnly)
	}
}

func TestWithPath(t *testing.T) {
	config := cookies.CookieConfig{}

	expectedPath := "/test"
	option := cookies.WithPath(expectedPath)
	option(&config)

	if config.Path != expectedPath {
		t.Errorf("WithPath did not set the path correctly. Expected: %s, Got: %s", expectedPath, config.Path)
	}
}

func TestWithSecure(t *testing.T) {
	config := cookies.CookieConfig{}

	option := cookies.WithSecure(true)
	option(&config)

	if !config.Secure {
		t.Errorf("WithSecure did not set the flag correctly. Expected: %v, Got: %v", true, config.Secure)
	}

	option = cookies.WithSecure(false)
	option(&config)

	if config.Secure {
		t.Errorf("WithSecure did not set the flag correctly. Expected: %v, Got: %v", false, config.Secure)
	}
}

func TestWithSameSite(t *testing.T) {
	config := cookies.CookieConfig{}

	expectedSameSite := http.SameSiteStrictMode
	option := cookies.WithSameSite(expectedSameSite)
	option(&config)

	if config.SameSite != expectedSameSite {
		t.Errorf("WithSameSite did not set the policy correctly. Expected: %v, Got: %v", expectedSameSite, config.SameSite)
	}
}

func TestNewCookieConfig(t *testing.T) {
	testCases := []struct {
		name           string
		cookieName     string
		options        []cookies.CookieOption
		expectedConfig cookies.CookieConfig
	}{
		{
			name:       "Default configuration",
			cookieName: "test-cookie",
			options:    []cookies.CookieOption{},
			expectedConfig: cookies.CookieConfig{
				Name:     "test-cookie",
				Path:     "/",
				MaxAge:   24 * time.Hour,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			},
		},
		{
			name:       "With custom options",
			cookieName: "custom-cookie",
			options: []cookies.CookieOption{
				cookies.WithValue("custom-value"),
				cookies.WithMaxAge(1 * time.Hour),
				cookies.WithPath("/custom"),
				cookies.WithHttpOnly(false),
				cookies.WithSecure(true),
				cookies.WithSameSite(http.SameSiteStrictMode),
			},
			expectedConfig: cookies.CookieConfig{
				Name:     "custom-cookie",
				Value:    "custom-value",
				Path:     "/custom",
				MaxAge:   1 * time.Hour,
				HttpOnly: false,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := cookies.NewCookieConfig(tc.cookieName, tc.options...)

			if config.Name != tc.expectedConfig.Name {
				t.Errorf("Incorrect name. Expected: %s, Got: %s", tc.expectedConfig.Name, config.Name)
			}
			if config.Value != tc.expectedConfig.Value {
				t.Errorf("Incorrect value. Expected: %s, Got: %s", tc.expectedConfig.Value, config.Value)
			}
			if config.Path != tc.expectedConfig.Path {
				t.Errorf("Incorrect path. Expected: %s, Got: %s", tc.expectedConfig.Path, config.Path)
			}
			if config.MaxAge != tc.expectedConfig.MaxAge {
				t.Errorf("Incorrect MaxAge. Expected: %v, Got: %v", tc.expectedConfig.MaxAge, config.MaxAge)
			}
			if config.HttpOnly != tc.expectedConfig.HttpOnly {
				t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", tc.expectedConfig.HttpOnly, config.HttpOnly)
			}
			if config.Secure != tc.expectedConfig.Secure {
				t.Errorf("Incorrect Secure. Expected: %v, Got: %v", tc.expectedConfig.Secure, config.Secure)
			}
			if config.SameSite != tc.expectedConfig.SameSite {
				t.Errorf("Incorrect SameSite. Expected: %v, Got: %v", tc.expectedConfig.SameSite, config.SameSite)
			}
		})
	}
}

func TestNewAuthCookieConfig(t *testing.T) {
	testCases := []struct {
		name           string
		token          string
		isProduction   bool
		options        []cookies.CookieOption
		expectedSecure bool
	}{
		{
			name:           "Development environment",
			token:          "test-token",
			isProduction:   false,
			options:        []cookies.CookieOption{},
			expectedSecure: false,
		},
		{
			name:           "Production environment",
			token:          "prod-token",
			isProduction:   true,
			options:        []cookies.CookieOption{},
			expectedSecure: true,
		},
		{
			name:           "Custom options",
			token:          "custom-token",
			isProduction:   false,
			options:        []cookies.CookieOption{cookies.WithSecure(true)},
			expectedSecure: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := cookies.NewAuthCookieConfig(tc.token, tc.isProduction, tc.options...)

			if config.Name != "token" {
				t.Errorf("Incorrect name. Expected: %s, Got: %s", "token", config.Name)
			}
			if config.Value != tc.token {
				t.Errorf("Incorrect value. Expected: %s, Got: %s", tc.token, config.Value)
			}
			if config.Path != "/" {
				t.Errorf("Incorrect path. Expected: %s, Got: %s", "/", config.Path)
			}
			if config.MaxAge != 12*time.Hour {
				t.Errorf("Incorrect MaxAge. Expected: %v, Got: %v", 12*time.Hour, config.MaxAge)
			}
			if !config.HttpOnly {
				t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", true, config.HttpOnly)
			}
			if config.Secure != tc.expectedSecure {
				t.Errorf("Incorrect Secure. Expected: %v, Got: %v", tc.expectedSecure, config.Secure)
			}
			if config.SameSite != http.SameSiteLaxMode {
				t.Errorf("Incorrect SameSite. Expected: %v, Got: %v", http.SameSiteLaxMode, config.SameSite)
			}
		})
	}
}

func TestSetCookie(t *testing.T) {
	w := httptest.NewRecorder()

	config := cookies.CookieConfig{
		Name:     "test-cookie",
		Value:    "test-value",
		MaxAge:   1 * time.Hour,
		HttpOnly: true,
		Path:     "/test",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	cookies.SetCookie(w, config)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Incorrect number of cookies. Expected: %d, Got: %d", 1, len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != config.Name {
		t.Errorf("Incorrect name. Expected: %s, Got: %s", config.Name, cookie.Name)
	}
	if cookie.Value != config.Value {
		t.Errorf("Incorrect value. Expected: %s, Got: %s", config.Value, cookie.Value)
	}
	if cookie.Path != config.Path {
		t.Errorf("Incorrect path. Expected: %s, Got: %s", config.Path, cookie.Path)
	}
	if cookie.MaxAge != 0 {
		t.Errorf("Incorrect MaxAge. Expected: %v, Got: %v", 0, cookie.MaxAge)
	}
	if !cookie.HttpOnly {
		t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", true, cookie.HttpOnly)
	}
	if !cookie.Secure {
		t.Errorf("Incorrect Secure. Expected: %v, Got: %v", true, cookie.Secure)
	}
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Errorf("Incorrect SameSite. Expected: %v, Got: %v", http.SameSiteStrictMode, cookie.SameSite)
	}

	expectedTime := time.Now().Add(config.MaxAge)
	timeDiff := cookie.Expires.Sub(expectedTime)
	if timeDiff < -1*time.Second || timeDiff > 1*time.Second {
		t.Errorf("Incorrect Expires. Expected around: %v, Got: %v", expectedTime, cookie.Expires)
	}
}

func TestSetAuthCookie(t *testing.T) {
	testCases := []struct {
		name         string
		token        string
		isProduction bool
		options      []cookies.CookieOption
	}{
		{
			name:         "Authentication cookie in development",
			token:        "test-token",
			isProduction: false,
			options:      []cookies.CookieOption{},
		},
		{
			name:         "Authentication cookie in production",
			token:        "prod-token",
			isProduction: true,
			options:      []cookies.CookieOption{},
		},
		{
			name:         "Authentication cookie with custom options",
			token:        "custom-token",
			isProduction: false,
			options:      []cookies.CookieOption{cookies.WithMaxAge(2 * time.Hour)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			cookies.SetAuthCookie(w, tc.token, tc.isProduction, tc.options...)

			cookies := w.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("Incorrect number of cookies. Expected: %d, Got: %d", 1, len(cookies))
			}

			cookie := cookies[0]
			if cookie.Name != "token" {
				t.Errorf("Incorrect name. Expected: %s, Got: %s", "token", cookie.Name)
			}
			if cookie.Value != tc.token {
				t.Errorf("Incorrect value. Expected: %s, Got: %s", tc.token, cookie.Value)
			}
			if cookie.Path != "/" {
				t.Errorf("Incorrect path. Expected: %s, Got: %s", "/", cookie.Path)
			}
			if !cookie.HttpOnly {
				t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", true, cookie.HttpOnly)
			}
		})
	}
}

func TestClearCookie(t *testing.T) {
	w := httptest.NewRecorder()

	config := cookies.CookieConfig{
		Name:     "test-cookie",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}

	cookies.SetCookie(w, config)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Incorrect number of cookies. Expected: %d, Got: %d", 1, len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != config.Name {
		t.Errorf("Incorrect name. Expected: %s, Got: %s", config.Name, cookie.Name)
	}

	if cookie.MaxAge != -1 {
		t.Errorf("Incorrect MaxAge for deletion. Expected: %v, Got: %v", -1, cookie.MaxAge)
	}

	expiredTime := time.Time{}
	if !cookie.Expires.Equal(expiredTime) && !cookie.Expires.Before(time.Now()) {
		t.Errorf("Incorrect Expires for deletion. Must be zero or in the past, Got: %v", cookie.Expires)
	}
}

func ExampleSetAuthCookie() {
	w := httptest.NewRecorder()

	cookies.SetAuthCookie(w, "user-token", false)
}

func ExampleNewCookieConfig_userPreferences() {
	w := httptest.NewRecorder()

	config := cookies.NewCookieConfig("preferences",
		cookies.WithValue("theme=dark"),
		cookies.WithMaxAge(30*24*time.Hour),
		cookies.WithHttpOnly(false),
	)

	cookies.SetCookie(w, config)
}

func ExampleNewCookieConfig_secureCookie() {
	w := httptest.NewRecorder()

	config := cookies.NewCookieConfig("secure-data",
		cookies.WithValue("sensitive-information"),
		cookies.WithMaxAge(1*time.Hour),
		cookies.WithHttpOnly(true),
		cookies.WithSecure(true),
		cookies.WithSameSite(http.SameSiteStrictMode),
	)

	cookies.SetCookie(w, config)
}
