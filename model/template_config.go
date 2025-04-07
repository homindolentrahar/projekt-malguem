package model

type TemplateConfig struct {
	Name      string                            `yaml:"name"`
	Language  string                            `yaml:"language"`
	FileCase  string                            `yaml:"file_case"`
	DirCase   string                            `yaml:"dir_case"`
	Variables map[string]TemplateConfigVariable `yaml:"variables"`
}

type TemplateConfigVariable struct {
	Type    any    `yaml:"type"`
	Prompt  string `yaml:"prompt"`
	Default any    `yaml:"default"`
}
