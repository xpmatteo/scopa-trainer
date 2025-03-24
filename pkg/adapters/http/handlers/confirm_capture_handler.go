package handlers

import (
	"net/http"
)

// ConfirmCaptureService is an interface for confirming a capture
type ConfirmCaptureService interface {
	ConfirmCapture()
}

// NewHandleConfirmCapture creates a new confirm capture handler
func NewHandleConfirmCapture(service ConfirmCaptureService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Confirm the capture
		service.ConfirmCapture()

		// Redirect back to the main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
