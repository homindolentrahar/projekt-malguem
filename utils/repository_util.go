package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Download a repository into a zip file with Github's Tarball feature
func DownloadZip(url, subdir, ref, zipPath, outputPath string) error {
	// Default branch is `main`
	if ref == "" {
		ref = "master"
	}

	os.MkdirAll("temps", os.ModePerm)

	// Download the tarball file
	tarballURL := fmt.Sprintf("%s/archive/refs/heads/%s.zip", url, ref)
	response, err := http.Get(tarballURL)
	HandleErrorExit(err)
	defer response.Body.Close()

	// Make sure the response is OK
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download %s: %s", tarballURL, response.Status)
	}

	tempZipPath := fmt.Sprintf("temps/%s", zipPath)
	zipFile, err := os.Create(tempZipPath)
	HandleErrorReturn(err)
	defer zipFile.Close()

	_, err = io.Copy(zipFile, response.Body)

	UnzipSubdir(tempZipPath, subdir, outputPath)

	os.Remove(tempZipPath)

	return err
}
