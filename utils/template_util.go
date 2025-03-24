package utils

import (
	"github.com/cbroglie/mustache"
	"os"
)

func renderTemplate(template string, data map[string]any) (string, error) {
	file, err := os.ReadFile(template)
	if err != nil {
		return "", err
	}

	result, err := mustache.Render(string(file), data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func ExecuteTemplate(source, target string, data map[string]any) error {
	result, err := renderTemplate(source, data)
	if err != nil {
		return err
	}

	return os.WriteFile(target, []byte(result), 0644)
}
