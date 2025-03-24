package handlers

import (
	"net/http"
)

// TestCombinationScenario is an interface for setting up a test scenario for combination captures
type TestCombinationScenario interface {
	SetupCombinationTest()
}

// NewHandleTestCombination creates a new test combination handler
func NewHandleTestCombination(service TestCombinationScenario) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up the combination test scenario
		service.SetupCombinationTest()

		// Redirect to the main game page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
