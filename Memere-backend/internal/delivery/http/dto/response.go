package dto

// Paginated is the standard envelope for cursor-paginated list endpoints
// (spec §5.4). NextCursor is empty on the last page. Data holds the page's
// items (already mapped to response DTOs).
type Paginated struct {
	Data       any    `json:"data"`
	NextCursor string `json:"next_cursor"`
	Limit      int    `json:"limit"`
}

// ErrorBody is the standard error envelope: { code, message, details }
// (spec §5.4). The delivery layer's respondError builds this from an AppError.
type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}
