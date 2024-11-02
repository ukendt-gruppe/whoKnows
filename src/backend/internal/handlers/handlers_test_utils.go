// File: src/backend/internal/handlers/handlers_test_utils.go

package handlers

import (
	"html/template"
	"testing"
)

func SetupTestTemplates(t *testing.T) {
	t.Helper()
	if GetTemplates() == nil {
		templates := template.Must(template.New("mock").Parse("{{.}}"))
		SetTemplates(templates)
	}
}
