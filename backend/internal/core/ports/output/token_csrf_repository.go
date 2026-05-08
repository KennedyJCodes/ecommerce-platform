// Package output defines output port interfaces for the application.
package output

import "github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"

// CSRFRepository defines the persistence contract for managing CSRF token lifecycle.
// Implementations handle the secure storage, retrieval, and cleanup of anti-forgery tokens.
type CSRFRepository interface {
	// Save persists a new CSRF token associated with a user or session.
    // Parameters:
    //   - token: Pointer to the CSRFToken model containing value and expiration data.
    // Returns:
    //   - error: non-nil if the storage operation fails.
	Save(token *models.CSRFToken) error

	// Find retrieves a specific CSRF token record by its value.
    // Parameters:
    //   - value: The unique string value of the token to search for.
    // Returns:
    //   - *models.CSRFToken: pointer to the found token model.
    //   - error: non-nil if the token does not exist or the query fails.
	Find(value string) (*models.CSRFToken, error)

	// Delete removes a CSRF token from the storage.
    // This should be used during logout or after a token has been consumed/expired.
    // Parameters:
    //   - value: The unique string value of the token to be deleted.
    // Returns:
    //   - error: non-nil if the deletion fails.
	Delete(value string) error
}
