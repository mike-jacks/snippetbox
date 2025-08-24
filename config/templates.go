package config

import "github.com/mike-jacks/snippetbox/internal/models"

// templateData is the holding structure for any dynamic data
// to pass to HTML temlates
type templateData struct {
	Snippet models.Snippet
}