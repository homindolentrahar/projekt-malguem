package utils

import (
	"os"
	"regexp"

	"github.com/cbroglie/mustache"
)

var sectionPattern = regexp.MustCompile(`{{#([\w]+)}}([\w]+){{/[\w]+}}`)

func renderTemplate(template string, data map[string]string) (string, error) {
	file, err := os.ReadFile(template)
	if err != nil {
		return "", err
	}

	templateString := string(file)
	preprocessedTemplate := sectionPattern.ReplaceAllStringFunc(templateString, func(match string) string {
		// Extract format type and variable name
		// {{format}} varName {{/format}}
		//     1         2         3
		matches := sectionPattern.FindStringSubmatch(match)
		// Invalid format
		if len(matches) < 3 {
			return match
		}

		format, variable := matches[1], matches[2]

		// Check if variable exists inside userData
		value, exists := data[variable]
		if !exists {
			return match
		}

		switch format {
		case "pascal_case":
			return ToPascalCase(value)
		case "camel_case":
			return ToCamelCase(value)
		case "snake_case":
			return ToSnakeCase(value)
		case "kebab_case":
			return ToKebabCase(value)
		default:
			return value
		}
	})

	result, err := mustache.Render(preprocessedTemplate, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func ExecuteTemplate(source, target string, data map[string]string) error {
	result, err := renderTemplate(source, data)
	if err != nil {
		return err
	}

	return os.WriteFile(target, []byte(result), 0644)
}
