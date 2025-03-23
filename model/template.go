package model

type Template struct {
	Name      string                      `yaml:"name"`
	Language  string                      `yaml:"language"`
	Variables map[string]TemplateVariable `yaml:"variables"`
}

type TemplateVariable struct {
	Type    any `yaml:"type"`
	Default any `yaml:"default"`
}
