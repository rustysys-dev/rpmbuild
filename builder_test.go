package rpmbuild_test

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/go-git/v5"
	"github.com/rustysys-dev/rpmbuild"
)

func TestGenRPMName(t *testing.T) {
}

func TestSetNameFromRepo(t *testing.T) {
	t.Run("Successful Open", func(t *testing.T) {
		tmpDir, fn := helpSetupGitRun(t)
		defer fn()

		builder := &rpmbuild.Builder{}
		err := builder.SetNameFromRepo()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expectedName := filepath.Base(tmpDir)
		if builder.Name != expectedName {
			t.Errorf("Incorrect repository name. Got: %s, Want: %s", builder.Name, expectedName)
		}
	})
}

func helpSetupRun(t *testing.T) (string, func()) {
	oldDir := os.Getenv("PWD")
	tmpDir, err := os.MkdirTemp("", "git-repo-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	return tmpDir, func() {
		os.RemoveAll(tmpDir)
		os.Chdir(oldDir)
	}
}

func helpSetupGitRun(t *testing.T) (string, func()) {
	tmpDir, fn := helpSetupRun(t)

	_, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to initialize Git repository: %v", err)
	}
	return tmpDir, fn
}

func recursiveCopy(oldDir, tmpDir string) error {
	return filepath.WalkDir(oldDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path '%s': %w", path, err)
		}

		relPath, err := filepath.Rel(oldDir, path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}
		newPath := filepath.Join(tmpDir, relPath)

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("error getting file info: %w", err)
		}

		if info.IsDir() {
			return os.MkdirAll(newPath, info.Mode())
		}

		src, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening source file: %w", err)
		}
		defer src.Close()

		dst, err := os.Create(newPath)
		if err != nil {
			return fmt.Errorf("error creating destination file: %w", err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		return err
	})
}
