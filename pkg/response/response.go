// Package response provides standard HTTP response envelopes for all
// services in the orion-v2 monorepo.
//
// Every endpoint returns one of two shapes:
//
//	Success  → { "success": true,  "message": "...", "data": <T>, "meta": <Meta|null> }
//	Error    → { "success": false, "message": "...", "error": { "code": "...", "details": <any|null> } }
package response

// ─── Success ─────────────────────────────────────────────────────────────────

// Response is the standard success envelope.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Meta carries pagination information for list endpoints.
type Meta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// NewMeta builds a Meta from raw pagination values.
func NewMeta(total, page, limit int) *Meta {
	if limit <= 0 {
		limit = 1
	}
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	return &Meta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// Success builds a standard success response with data.
func Success(message string, data any) Response {
	return Response{Success: true, Message: message, Data: data}
}

// SuccessWithMeta builds a success response with pagination metadata.
func SuccessWithMeta(message string, data any, meta *Meta) Response {
	return Response{Success: true, Message: message, Data: data, Meta: meta}
}

// Created builds a 201-appropriate success response.
func Created(data any) Response {
	return Response{Success: true, Message: "created", Data: data}
}

// ─── Error ────────────────────────────────────────────────────────────────────

// ErrorResponse is the standard error envelope.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail carries a machine-readable code and optional validation details.
type ErrorDetail struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

// Error code constants — use these in handlers for consistent client-side handling.
const (
	CodeBadRequest   = "BAD_REQUEST"
	CodeValidation   = "VALIDATION_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
)

// Fail builds a generic error response.
func Fail(message, code string, details any) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
		Error:   ErrorDetail{Code: code, Details: details},
	}
}

// ErrBadRequest returns a 400 error response.
func ErrBadRequest(message string) ErrorResponse {
	return Fail(message, CodeBadRequest, nil)
}

// ErrValidation returns a 422 error response with optional field-level details.
func ErrValidation(details any) ErrorResponse {
	return Fail("validation failed", CodeValidation, details)
}

// ErrUnauthorized returns a 401 error response.
func ErrUnauthorized() ErrorResponse {
	return Fail("unauthorized", CodeUnauthorized, nil)
}

// ErrForbidden returns a 403 error response.
func ErrForbidden() ErrorResponse {
	return Fail("forbidden", CodeForbidden, nil)
}

// ErrNotFound returns a 404 error response.
func ErrNotFound(resource string) ErrorResponse {
	return Fail(resource+" not found", CodeNotFound, nil)
}

// ErrConflict returns a 409 error response.
func ErrConflict(message string) ErrorResponse {
	return Fail(message, CodeConflict, nil)
}

// ErrInternal returns a 500 error response.
// Never leak internal error details — log them server-side instead.
func ErrInternal() ErrorResponse {
	return Fail("an unexpected error occurred", CodeInternal, nil)
}
