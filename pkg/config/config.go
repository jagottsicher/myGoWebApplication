package config

import "html/template"

// AppConfig is a struct holding die application's configuration
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
}
