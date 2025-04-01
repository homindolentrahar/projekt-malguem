package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
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

func UnzipFile(source, target string) error {
	zipReader, err := zip.OpenReader(source)
	HandleErrorExit(err)
	defer zipReader.Close()

	os.MkdirAll(target, os.ModePerm)

	for _, file := range zipReader.File {
		targetPath := filepath.Join(target, file.Name)

		// If the entry is folder, then create
		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, os.ModePerm)
			continue
		}

		err := extractZipFile(file, targetPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnzipSubdir(source, subdir, target string) error {
	zipReader, err := zip.OpenReader(source)
	HandleErrorExit(err)
	defer zipReader.Close()

	os.MkdirAll(target, os.ModePerm)

	// Normalize the subdir path
	subdir = filepath.ToSlash(subdir)

	// Add the trailing slash if not exist
	if !strings.HasSuffix(subdir, "/") {
		subdir += "/"
	}

	// Find the root folder
	var rootFolder string
	for _, file := range zipReader.File {
		normalizedPath := filepath.ToSlash(file.Name)
		parts := strings.SplitN(normalizedPath, "/", 2)

		if len(parts) > 1 {
			rootFolder = parts[0]
			break
		}
	}
	log.Debug().Msgf("Root: %s", rootFolder)

	// Extract the zip file
	targetDir := fmt.Sprintf("%s/%s", rootFolder, subdir)
	log.Debug().Msgf("Target: %s", targetDir)
	for _, file := range zipReader.File {
		normalizedPath := filepath.ToSlash(file.Name)

		if !strings.HasPrefix(normalizedPath, targetDir) {
			continue
		}

		relativePath := strings.Trim(normalizedPath, targetDir)
		targetPath := filepath.Join(target, subdir, filepath.FromSlash(relativePath))

		// if the entry is folder, then make sure the folder is available
		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, os.ModePerm)
			continue
		}

		err := extractZipFile(file, targetPath)
		HandleErrorReturn(err)
	}

	return nil
}

func extractZipFile(source *zip.File, target string) error {
	os.MkdirAll(filepath.Dir(target), os.ModePerm)

	srcFile, err := source.Open()
	HandleErrorExit(err)
	defer srcFile.Close()

	destFile, err := os.Create(target)
	HandleErrorExit(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	return err
}
