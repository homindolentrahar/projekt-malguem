package github

import (
	"fmt"
	"malguem/utils"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneRepo(url, target, subdir string) error {
	response, err := http.Get(url)
	utils.HandleErrorReturn(err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch the repository: %s", response.Status)
	}

	// Make sure the target directory exists
	os.MkdirAll(target, os.ModePerm)

	// Extract the subdirectory from the tarball
	return utils.ExtractSubdirectory(response.Body, target, subdir)
}

func CloneSubdir(url, branch, path, output string) error {
	tempDir := "temp"

	// Clone the repo without checking out files
	cmd := exec.Command("git", "clone", "--no-checkout", "--depth", "1", "--filter=blob:none", "--branch", branch, url, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	// Enable sparse checkout
	cmd = exec.Command("git", "-C", tempDir, "sparse-checkout", "init", "--cone")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable sparse checkout: %w", err)
	}

	// Speciffy the subdirectory
	cmd = exec.Command("git", "-C", tempDir, "sparse-checkout", "set", path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to specify the subdirectory: %w", err)
	}

	// Checkout only the specified folder
	cmd = exec.Command("git", "-C", tempDir, "checkout")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout the specified folder: %w", err)
	}

	// Move the folder to output
	sourcePath := filepath.Join(tempDir, path)
	if err := utils.CopyDir(sourcePath, output); err != nil {
		return fmt.Errorf("failed to move the folder: %w", err)
	}

	// Cleanup
	os.RemoveAll(tempDir)

	return nil
}
