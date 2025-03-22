package utils

import (
	"os"
	"path/filepath"
)

func GetCurrentDir() (string, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Base(workingDirectory), nil
}