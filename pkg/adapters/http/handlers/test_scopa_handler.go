package handlers

import (
	"net/http"
)

// TestScopaScenario is an interface for setting up a test scenario for scopa
type TestScopaScenario interface {
	SetupScopaTestScenario()
}

// NewHandleTestScopa creates a new test scopa handler
func NewHandleTestScopa(service TestScopaScenario) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up the scopa test scenario
		service.SetupScopaTestScenario()

		// Redirect to the main game page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
