// Package modelsDatabse contains models stored temporarily in Redis.
package modelsdb

// PendingUser represents a user registration awaiting email verification.
// Redis should store this model with a TTL derived from ExpiresAt.
type PendingUser struct {
	ID           string `redis:"id"`
	Username     string `redis:"username"`
	PasswordHash string `redis:"password_hash"`
	Email        string `redis:"email"`
	HashCode     string `redis:"hash_code"`
	Attempts     int    `redis:"attempts"`
}
