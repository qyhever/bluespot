package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"bluespot/internal/config"
)

func TestUploadRepositoryDeleteFinalFile(t *testing.T) {
	repo := NewUploadRepository()
	dirs := newTestUploadRepoDirs(t)

	if err := os.MkdirAll(dirs.large, 0755); err != nil {
		t.Fatal(err)
	}
	finalPath := filepath.Join(dirs.large, "demo.txt")
	if err := os.WriteFile(finalPath, []byte("demo"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := repo.DeleteFinalFile("demo.txt"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(finalPath); !os.IsNotExist(err) {
		t.Fatalf("expected final file deleted, stat err = %v", err)
	}
}

func TestUploadRepositoryDeleteFinalFileRejectsPathTraversal(t *testing.T) {
	repo := NewUploadRepository()
	dirs := newTestUploadRepoDirs(t)

	outsidePath := filepath.Join(dirs.root, "outside.txt")
	if err := os.WriteFile(outsidePath, []byte("keep"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := repo.DeleteFinalFile("../outside.txt"); err == nil {
		t.Fatalf("expected invalid file name error")
	}
	if _, err := os.Stat(outsidePath); err != nil {
		t.Fatalf("expected outside file to remain, stat err = %v", err)
	}
}

type testUploadRepoDirs struct {
	root  string
	large string
}

func newTestUploadRepoDirs(t *testing.T) testUploadRepoDirs {
	t.Helper()
	tmp := t.TempDir()
	dirs := testUploadRepoDirs{
		root:  tmp,
		large: filepath.Join(tmp, "larges"),
	}
	oldConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		Attach: config.AttachConfig{
			UploadLargeFilePath: dirs.large,
			ChunkDirPath:        filepath.Join(tmp, "chunks"),
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = oldConfig
	})
	return dirs
}
