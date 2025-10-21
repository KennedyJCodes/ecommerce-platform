// Package models defines core domain entities for the sale-watches application.
// It contains simple data structures that represent business objects throughout the system.
package models

// Account represents user credentials used for authentication and registration.

// Fields:
//   - UserName: the unique identifier chosen by the user for login purposes.
//   - Password: the user's secret passphrase; it should be transmitted securely (e.g., over HTTPS) and never logged or stored in plain text. In the domain, it is combined with a salt and hashed before persistence.
type Account struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}
