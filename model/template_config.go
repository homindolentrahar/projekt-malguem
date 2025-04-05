package model

type TemplateConfig struct {
	Name      string                            `yaml:"name"`
	Language  string                            `yaml:"language"`
	Variables map[string]TemplateConfigVariable `yaml:"variables"`
}

type TemplateConfigVariable struct {
	Type    any    `yaml:"type"`
	Prompt  string `yaml:"prompt"`
	Default any    `yaml:"default"`
}
