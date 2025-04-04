package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractSubdirectory(reader io.ReadCloser, target, subDir string) error {
	gzReader, err := gzip.NewReader(reader)
	HandleErrorReturn(err)
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		// End of archive
		if err == io.EOF {
			break
		}
		HandleErrorReturn(err)

		// Check if file is in the subdirectory

		if strings.Contains(header.Name, subDir) {
			rootDir := strings.Split(header.Name, "/")[0]
			rootPath := filepath.Join(rootDir, subDir)
			relativePath := strings.TrimPrefix(header.Name, rootPath)
			targetPath := filepath.Join(target, relativePath)

			// Create directories if needed
			if header.Typeflag == tar.TypeDir {
				os.MkdirAll(targetPath, os.ModePerm)
				continue
			}

			// Create the target file and write the content
			targetFile, err := os.Create(targetPath)
			HandleErrorReturn(err)
			defer targetFile.Close()

			_, err = io.Copy(targetFile, tarReader)
			HandleErrorReturn(err)
		}
	}

	return nil
}
