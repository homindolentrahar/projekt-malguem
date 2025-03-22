package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CreateZipFile(source, target string) error {
	file, err := os.Create(target)
	HandleErrorExit(err)
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create zip entry
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Open file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Write file to zip
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)
		return err
	})

	return err
}
