package http

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	stdhttp "net/http"

	appErrors "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
)

const (
	// MaxAuthBodySize caps login/register JSON payloads to the small data they need.
	MaxAuthBodySize int64 = 8 * 1024

	// MaxCommentBodySize allows longer comment text while still rejecting oversized bodies.
	MaxCommentBodySize int64 = 16 * 1024
)

// DecodeJSONBody validates and decodes a single JSON request body into dst.
//
// The helper centralizes API body rules: Content-Type must be application/json,
// the stream is bounded with MaxBytesReader, unknown fields are rejected, and
// requests containing more than one JSON value are treated as invalid.
func DecodeJSONBody(w stdhttp.ResponseWriter, r *stdhttp.Request, dst any, maxBytes int64) error {
	if !hasJSONContentType(r) {
		return appErrors.NewUnsupportedMediaTypeError(appErrors.ErrUnsupportedMediaType)
	}

	// MaxBytesReader stops oversized payloads before the decoder consumes unbounded data.
	r.Body = stdhttp.MaxBytesReader(w, r.Body, maxBytes)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return decodeJSONError(err)
	}

	// A second successful decode means the client sent multiple JSON values.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return appErrors.NewBadRequestError(appErrors.ErrInvalidRequest)
	}

	return nil
}

func hasJSONContentType(r *stdhttp.Request) bool {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	return err == nil && mediaType == "application/json"
}

func decodeJSONError(err error) error {
	var maxBytesErr *stdhttp.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return appErrors.NewRequestEntityTooLargeError(appErrors.ErrRequestBodyTooLarge)
	}

	return appErrors.NewBadRequestError(appErrors.ErrInvalidRequest)
}
