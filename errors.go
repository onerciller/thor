package thor

import "net/http"

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(status int, message ...interface{}) *HTTPError {
	he := &HTTPError{Status: status, Message: http.StatusText(status)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

// Error returns the error message.
func (e *HTTPError) Error() interface{} {
	return e.Message
}

// StatusCode returns the HTTP status code.
func (e *HTTPError) StatusCode() int {
	return e.Status
}
