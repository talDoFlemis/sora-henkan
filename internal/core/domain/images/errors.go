package images

import "fmt"

// NonRetryableError represents an error that should not be retried
// When this error is returned, the message should be ACKed instead of NACKed
type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string {
	return fmt.Sprintf("non-retryable error: %v", e.Err)
}

func (e *NonRetryableError) Unwrap() error {
	return e.Err
}

// NewNonRetryableError creates a new NonRetryableError
func NewNonRetryableError(err error) *NonRetryableError {
	return &NonRetryableError{Err: err}
}
