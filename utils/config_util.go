package utils

import (
	"gopkg.in/yaml.v3"
	"malguem/model"
	"os"
)

func ReadMalguemConfig(configFile ...string) (*model.MalguemConfig, error) {
	config := "malguem.yaml"

	if len(configFile) > 0 {
		config = configFile[0]
	}

	file, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var malguemConfig model.MalguemConfig
	err = yaml.Unmarshal(file, &malguemConfig)
	if err != nil {
		return nil, err
	}

	return &malguemConfig, nil
}

func ReadTemplate(path string) (*model.Template, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var template model.Template
	err = yaml.Unmarshal(file, &template)

	if err != nil {
		return nil, err
	}

	return &template, err
}
