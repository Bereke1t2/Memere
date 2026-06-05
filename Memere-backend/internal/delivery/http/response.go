package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// respondJSON writes a success payload with the given status.
func respondJSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

// respondError maps an error to the standard { code, message, details } envelope
// (spec §5.4). A typed *apperror.AppError carries its own HTTP status, code,
// message and details. Any other error is logged and collapsed to a generic 500
// so internal details never reach the client (spec §7.3).
func respondError(c *gin.Context, err error) {
	var ae *apperror.AppError
	if errors.As(err, &ae) {
		c.JSON(ae.HTTPStatus, dto.ErrorBody{
			Code:    ae.Code,
			Message: ae.Message,
			Details: ae.Details,
		})
		return
	}

	// Unexpected: log the real error server-side, return an opaque 500.
	log.Printf("unhandled error: %v", err)
	c.JSON(http.StatusInternalServerError, dto.ErrorBody{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An internal server error occurred",
	})
}
