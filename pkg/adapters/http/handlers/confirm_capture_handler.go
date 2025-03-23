package handlers

import (
	"net/http"
)

// ConfirmCaptureService is an interface for confirming a capture
type ConfirmCaptureService interface {
	ConfirmCapture()
}

// HandleConfirmCapture handles confirming a capture selection
type HandleConfirmCapture struct {
	service ConfirmCaptureService
}

// NewHandleConfirmCapture creates a new confirm capture handler
func NewHandleConfirmCapture(service ConfirmCaptureService) http.HandlerFunc {
	handler := &HandleConfirmCapture{
		service: service,
	}
	return handler.ServeHTTP
}

// ServeHTTP handles the HTTP request
func (h *HandleConfirmCapture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Confirm the capture
	h.service.ConfirmCapture()

	// Redirect back to the main page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}