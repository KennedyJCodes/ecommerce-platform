// Package input defines service contracts for comment-related business logic, user operations, and input validation.
package input

// CSRFService defines the contract for Cross-Site Request Forgery protection.
// Implementations must ensure the integrity and authenticity of requests by linking tokens to specific users.
type CSRFService interface {
	// GenerateToken creates a unique CSRF token for a specific user session.
    // Parameters:
    //   - userID: The unique identifier of the user to bind the token to.
    // Returns:
    //   - string: A secure, random token value.
    //   - error: non-nil if the token generation fails.
	GenerateToken(userID string) (string, error)

	// ValidateToken verifies if a provided token is valid and belongs to the specified user.
    // Parameters:
    //   - tokenValue: The token string received from the client request.
    //   - userID: The identifier of the user performing the action.
    // Returns:
    //   - error: non-nil if the token is expired, invalid, or does not match the user.
	ValidateToken(tokenValue string, userID string) error
}
