package utils

import (
	"io"
	"io/fs"
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

func CopyDir(source, target string) error {
	// Ensure the target directory exists
	err := os.MkdirAll(target, os.ModePerm)
	HandleErrorReturn(err)

	return filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		HandleErrorReturn(err)

		relativePath, err := filepath.Rel(source, path)
		HandleErrorReturn(err)

		targetPath := filepath.Join(target, relativePath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return CopyFile(path, targetPath)
	})
}

func CopyFile(source, target string) error {
	sourceFile, err := os.Open(source)
	HandleErrorReturn(err)
	defer sourceFile.Close()

	targetFile, err := os.Create(target)
	HandleErrorReturn(err)
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)

	return err
}
