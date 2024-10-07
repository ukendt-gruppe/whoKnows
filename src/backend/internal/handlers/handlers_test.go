// File: src/backend/internal/handlers/handlers_test.go

package handlers

import (
    "html/template"
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // Override the template initialization for tests
    templates = template.Must(template.New("mock").Parse("{{.}}"))
    
    // Run the tests
    code := m.Run()

    // Exit with the test result code
    os.Exit(code)
}