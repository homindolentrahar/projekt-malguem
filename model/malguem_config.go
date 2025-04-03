package model

type MalguemConfig struct {
	Name      string                     `yaml:"name"`
	Templates map[string]MalguemTemplate `yaml:"templates"`
}

type MalguemTemplate struct {
	Path   *string       `yaml:"path,omitempty"`
	Github *GithubSource `yaml:"github,omitempty"`
	Output string        `yaml:"output,omitempty"`
}

type GithubSource struct {
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
	Ref  string `yaml:"ref,omitempty"`
}
