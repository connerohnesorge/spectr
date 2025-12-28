package initialize

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//nolint:revive // cognitive-complexity - comprehensive test coverage
//nolint:revive // cognitive-complexity - comprehensive test
func TestExpandPath(t *testing.T) {
	t.Run(
		"expands home directory path",
		func(t *testing.T) {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				t.Fatalf(
					"failed to get home dir: %v",
					err,
				)
			}

			expanded, err := ExpandPath("~/test")
			if err != nil {
				t.Fatalf(
					"ExpandPath failed: %v",
					err,
				)
			}

			expected := filepath.Join(
				homeDir,
				"test",
			)
			if expanded != expected {
				t.Errorf(
					"expected %s, got %s",
					expected,
					expanded,
				)
			}
		},
	)

	t.Run(
		"expands tilde alone to home directory",
		func(t *testing.T) {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				t.Fatalf(
					"failed to get home dir: %v",
					err,
				)
			}

			expanded, err := ExpandPath("~")
			if err != nil {
				t.Fatalf(
					"ExpandPath failed: %v",
					err,
				)
			}

			if expanded != homeDir {
				t.Errorf(
					"expected %s, got %s",
					homeDir,
					expanded,
				)
			}
		},
	)

	t.Run(
		"converts relative path to absolute",
		func(t *testing.T) {
			expanded, err := ExpandPath("./test")
			if err != nil {
				t.Fatalf(
					"ExpandPath failed: %v",
					err,
				)
			}

			if !filepath.IsAbs(expanded) {
				t.Errorf(
					"expected absolute path, got %s",
					expanded,
				)
			}

			if !strings.HasSuffix(
				expanded,
				"test",
			) {
				t.Errorf(
					"expected path ending with 'test', got %s",
					expanded,
				)
			}
		},
	)

	t.Run(
		"returns absolute path unchanged",
		func(t *testing.T) {
			absPath := "/absolute/path/test"
			expanded, err := ExpandPath(absPath)
			if err != nil {
				t.Fatalf(
					"ExpandPath failed: %v",
					err,
				)
			}

			if expanded != absPath {
				t.Errorf(
					"expected %s, got %s",
					absPath,
					expanded,
				)
			}
		},
	)

	t.Run(
		"returns error for empty path",
		func(t *testing.T) {
			_, err := ExpandPath("")
			if err == nil {
				t.Error(
					"expected error for empty path",
				)
			}

			if !strings.Contains(
				err.Error(),
				"cannot be empty",
			) {
				t.Errorf(
					"unexpected error message: %v",
					err,
				)
			}
		},
	)
}

//nolint:revive // cognitive-complexity - comprehensive test coverage

//nolint:revive // cognitive-complexity - comprehensive test
func TestEnsureDir(t *testing.T) {
	t.Run(
		"creates directory",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testDir := filepath.Join(
				tmpDir,
				"testdir",
			)

			err := EnsureDir(testDir)
			if err != nil {
				t.Fatalf(
					"EnsureDir failed: %v",
					err,
				)
			}

			// Verify directory exists
			info, err := os.Stat(testDir)
			if err != nil {
				t.Fatalf(
					"directory was not created: %v",
					err,
				)
			}

			if !info.IsDir() {
				t.Error(
					"path exists but is not a directory",
				)
			}
		},
	)

	t.Run(
		"creates parent directories",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testDir := filepath.Join(
				tmpDir,
				"parent",
				"child",
				"grandchild",
			)

			err := EnsureDir(testDir)
			if err != nil {
				t.Fatalf(
					"EnsureDir failed: %v",
					err,
				)
			}

			// Verify all directories exist
			info, err := os.Stat(testDir)
			if err != nil {
				t.Fatalf(
					"directory was not created: %v",
					err,
				)
			}

			if !info.IsDir() {
				t.Error(
					"path exists but is not a directory",
				)
			}
		},
	)

	t.Run(
		"is idempotent - no error if directory exists",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testDir := filepath.Join(
				tmpDir,
				"existing",
			)

			// Create directory first time
			err := EnsureDir(testDir)
			if err != nil {
				t.Fatalf(
					"first EnsureDir failed: %v",
					err,
				)
			}

			// Create again - should not error
			err = EnsureDir(testDir)
			if err != nil {
				t.Fatalf(
					"second EnsureDir failed: %v",
					err,
				)
			}
		},
	)

	t.Run(
		"returns error for empty path",
		func(t *testing.T) {
			err := EnsureDir("")
			if err == nil {
				t.Error(
					"expected error for empty path",
				)
			}

			if !strings.Contains(
				err.Error(),
				"cannot be empty",
			) {
				t.Errorf(
					"unexpected error message: %v",
					err,
				)
			}
		},
	)
	//nolint:revive // cognitive-complexity - comprehensive test coverage
}

//nolint:revive // cognitive-complexity - comprehensive test
func TestWriteFile(t *testing.T) {
	t.Run(
		"creates file with content",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(
				tmpDir,
				"test.txt",
			)
			content := []byte("test content")

			err := WriteFile(testFile, content)
			if err != nil {
				t.Fatalf(
					"WriteFile failed: %v",
					err,
				)
			}

			// Verify file exists and has correct content
			readContent, err := os.ReadFile(
				testFile,
			)
			if err != nil {
				t.Fatalf(
					"failed to read created file: %v",
					err,
				)
			}

			if !bytes.Equal(readContent, content) {
				t.Errorf(
					"expected content %q, got %q",
					content,
					readContent,
				)
			}
		},
	)

	t.Run(
		"creates parent directories if needed",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(
				tmpDir,
				"parent",
				"child",
				"test.txt",
			)
			content := []byte("test content")

			err := WriteFile(testFile, content)
			if err != nil {
				t.Fatalf(
					"WriteFile failed: %v",
					err,
				)
			}

			// Verify file exists
			if !FileExists(testFile) {
				t.Error("file was not created")
			}
		},
	)

	t.Run(
		"returns error if file already exists",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(
				tmpDir,
				"existing.txt",
			)

			// Create file first time
			err := WriteFile(
				testFile,
				[]byte("first"),
			)
			if err != nil {
				t.Fatalf(
					"first WriteFile failed: %v",
					err,
				)
			}

			// Try to write again - should error
			err = WriteFile(
				testFile,
				[]byte("second"),
			)
			if err == nil {
				t.Error(
					"expected error when file already exists",
				)
			}

			if !strings.Contains(
				err.Error(),
				"already exists",
			) {
				t.Errorf(
					"unexpected error message: %v",
					err,
				)
			}
		},
	)

	t.Run(
		"returns error for empty path",
		func(t *testing.T) {
			err := WriteFile(
				"",
				[]byte("content"),
			)
			if err == nil {
				t.Error(
					"expected error for empty path",
				)
			}

			if !strings.Contains(
				err.Error(),
				"cannot be empty",
			) {
				t.Errorf(
					"unexpected error message: %v",
					err,
				)
			}
		},
	)
}

func TestFileExists(t *testing.T) {
	t.Run(
		"returns true for existing file",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(
				tmpDir,
				"test.txt",
			)

			// Create file
			err := os.WriteFile(
				testFile,
				[]byte("test"),
				0o644,
			)
			if err != nil {
				t.Fatalf(
					"failed to create test file: %v",
					err,
				)
			}

			if !FileExists(testFile) {
				t.Error(
					"FileExists returned false for existing file",
				)
			}
		},
	)

	t.Run(
		"returns true for existing directory",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			testDir := filepath.Join(
				tmpDir,
				"testdir",
			)

			// Create directory
			err := os.Mkdir(testDir, 0o755)
			if err != nil {
				t.Fatalf(
					"failed to create test directory: %v",
					err,
				)
			}

			if !FileExists(testDir) {
				t.Error(
					"FileExists returned false for existing directory",
				)
			}
		},
	)

	t.Run(
		"returns false for non-existent path",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			nonExistent := filepath.Join(
				tmpDir,
				"does-not-exist",
			)

			if FileExists(nonExistent) {
				t.Error(
					"FileExists returned true for non-existent path",
				)
			}
		},
	)

	t.Run(
		"returns false for empty path",
		func(t *testing.T) {
			if FileExists("") {
				t.Error(
					"FileExists returned true for empty path",
				)
			}
		},
	)
}

func TestIsSpectrInitialized(t *testing.T) {
	t.Run(
		"returns true when spectr/project.md exists",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			spectrDir := filepath.Join(
				tmpDir,
				"spectr",
			)
			projectFile := filepath.Join(
				spectrDir,
				"project.md",
			)

			// Create spectr directory and project.md
			err := os.Mkdir(spectrDir, 0o755)
			if err != nil {
				t.Fatalf(
					"failed to create spectr directory: %v",
					err,
				)
			}

			err = os.WriteFile(
				projectFile,
				[]byte("# Project"),
				0o644,
			)
			if err != nil {
				t.Fatalf(
					"failed to create project.md: %v",
					err,
				)
			}

			if !IsSpectrInitialized(tmpDir) {
				t.Error(
					"IsSpectrInitialized returned false for initialized project",
				)
			}
		},
	)

	t.Run(
		"returns false when spectr directory does not exist",
		func(t *testing.T) {
			tmpDir := t.TempDir()

			if IsSpectrInitialized(tmpDir) {
				t.Error(
					"IsSpectrInitialized returned true for uninitialized project",
				)
			}
		},
	)

	t.Run(
		"returns false when spectr exists but project.md does not",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			spectrDir := filepath.Join(
				tmpDir,
				"spectr",
			)

			// Create spectr directory but not project.md
			err := os.Mkdir(spectrDir, 0o755)
			if err != nil {
				t.Fatalf(
					"failed to create spectr directory: %v",
					err,
				)
			}

			if IsSpectrInitialized(tmpDir) {
				t.Error(
					"IsSpectrInitialized returned true without project.md",
				)
			}
		},
	)

	t.Run(
		"returns false for empty path",
		func(t *testing.T) {
			if IsSpectrInitialized("") {
				//nolint:revive // cognitive-complexity - comprehensive test coverage
				t.Error(
					"IsSpectrInitialized returned true for empty path",
				)
			}
		},
	)
}

//nolint:revive // cognitive-complexity - comprehensive test
func TestBackupFile(t *testing.T) {
	t.Run(
		"creates backup of existing file",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			originalFile := filepath.Join(
				tmpDir,
				"original.txt",
			)
			content := []byte("original content")

			// Create original file
			err := os.WriteFile(
				originalFile,
				content,
				0o644,
			)
			if err != nil {
				t.Fatalf(
					"failed to create original file: %v",
					err,
				)
			}

			// Create backup
			err = BackupFile(originalFile)
			if err != nil {
				t.Fatalf(
					"BackupFile failed: %v",
					err,
				)
			}

			// Find backup file (has timestamp suffix)
			entries, err := os.ReadDir(tmpDir)
			if err != nil {
				t.Fatalf(
					"failed to read directory: %v",
					err,
				)
			}

			var backupFound bool
			var backupPath string
			for _, entry := range entries {
				name := entry.Name()
				if !strings.HasPrefix(
					name,
					"original.txt.backup.",
				) {
					continue
				}
				backupFound = true
				backupPath = filepath.Join(
					tmpDir,
					name,
				)

				break
			}

			if !backupFound {
				t.Error(
					"backup file was not created",
				)

				return
			}

			// Verify backup has same content as original
			if backupPath == "" {
				return
			}
			backupContent, err := os.ReadFile(
				backupPath,
			)
			if err != nil {
				t.Fatalf(
					"failed to read backup file: %v",
					err,
				)
			}

			if !bytes.Equal(backupContent, content) {
				t.Errorf(
					"backup content %q does not match original %q",
					backupContent,
					content,
				)
			}
		},
	)

	t.Run(
		"does not error if file does not exist",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			nonExistent := filepath.Join(
				tmpDir,
				"does-not-exist.txt",
			)

			err := BackupFile(nonExistent)
			if err != nil {
				t.Errorf(
					"BackupFile should not error for non-existent file: %v",
					err,
				)
			}
		},
	)

	t.Run(
		"returns error for empty path",
		func(t *testing.T) {
			err := BackupFile("")
			if err == nil {
				t.Error(
					"expected error for empty path",
				)
			}

			if !strings.Contains(
				err.Error(),
				"cannot be empty",
			) {
				t.Errorf(
					"unexpected error message: %v",
					err,
				)
			}
		},
	)

	t.Run(
		"creates multiple backups with different timestamps",
		func(t *testing.T) {
			tmpDir := t.TempDir()
			originalFile := filepath.Join(
				tmpDir,
				"original.txt",
			)
			content := []byte("original content")

			// Create original file
			err := os.WriteFile(
				originalFile,
				content,
				0o644,
			)
			if err != nil {
				t.Fatalf(
					"failed to create original file: %v",
					err,
				)
			}

			// Create first backup
			err = BackupFile(originalFile)
			if err != nil {
				t.Fatalf(
					"first BackupFile failed: %v",
					err,
				)
			}

			// Create second backup
			err = BackupFile(originalFile)
			if err != nil {
				t.Fatalf(
					"second BackupFile failed: %v",
					err,
				)
			}

			// Count backup files
			entries, err := os.ReadDir(tmpDir)
			if err != nil {
				t.Fatalf(
					"failed to read directory: %v",
					err,
				)
			}

			backupCount := 0
			for _, entry := range entries {
				if strings.HasPrefix(
					entry.Name(),
					"original.txt.backup.",
				) {
					backupCount++
				}
			}

			if backupCount != 2 {
				t.Errorf(
					"expected 2 backup files, found %d",
					backupCount,
				)
			}
		},
	)
}
