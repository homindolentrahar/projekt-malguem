package command

type {{#pascal_case}}name{{/pascal_case}} struct {
	Server string
	Port   int
}

func New{{#pascal_case}}name{{/pascal_case}}(server string, port int) *{{#pascal_case}}name{{/pascal_case}} {
	// {{#snake_case}}name{{/snake_case}}
	// {{#kebab_case}}name{{/kebab_case}}
	// {{#camel_case}}name{{/camel_case}}
	return &{{#pascal_case}}name{{/pascal_case}}{
		Server: server,
		Port:   port,
	}
}
