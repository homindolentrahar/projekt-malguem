package model

type MalguemConfig struct {
	Name      string                     `yaml:"name"`
	Templates map[string]MalguemTemplate `yaml:"templates"`
}

type MalguemTemplate struct {
	Path string `yaml:"path,omitempty"`
	Url  string `yaml:"url,omitempty"`
}
