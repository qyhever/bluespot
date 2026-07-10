package service

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"bluespot/internal/config"
	"bluespot/internal/model"
	"bluespot/internal/repository/persistence"
)

func TestGenerateUploadIDStable(t *testing.T) {
	got1 := GenerateUploadID("abc", 1, "salt")
	got2 := GenerateUploadID("abc", 1, "salt")
	if got1 != got2 {
		t.Fatalf("expected stable upload id, got %q and %q", got1, got2)
	}
}

func TestGenerateUploadIDDifferentUser(t *testing.T) {
	got1 := GenerateUploadID("abc", 1, "salt")
	got2 := GenerateUploadID("abc", 2, "salt")
	if got1 == got2 {
		t.Fatalf("expected different upload id for different users")
	}
}

func TestUploadVerifyDetectsExistingFinalFile(t *testing.T) {
	svc, dirs := newTestUploadService(t)
	fileMD5 := "0123456789abcdef0123456789abcdef"
	if err := os.MkdirAll(dirs.large, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirs.large, fileMD5+".mp4"), []byte("done"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := svc.Verify(7, model.UploadVerifyRequest{
		FileMD5:  fileMD5,
		FileName: "demo.mp4",
		FileSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !got.IsExists {
		t.Fatalf("expected final file to exist")
	}
	if got.URL != "http://test.local/larges/0123456789abcdef0123456789abcdef.mp4" {
		t.Fatalf("unexpected url: %s", got.URL)
	}
}

func TestUploadVerifyScansAndSortsUploadedChunks(t *testing.T) {
	svc, dirs := newTestUploadService(t)
	fileMD5 := "0123456789abcdef0123456789abcdef"
	uploadID := GenerateUploadID(fileMD5, 7, "salt")
	chunkDir := filepath.Join(dirs.chunks, uploadID)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"2", "0", "note", "1"} {
		if err := os.WriteFile(filepath.Join(chunkDir, name), []byte(name), 0644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := svc.Verify(7, model.UploadVerifyRequest{
		FileMD5:  fileMD5,
		FileName: "demo.mp4",
		FileSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int{0, 1, 2}
	if !reflect.DeepEqual(got.UploadedChunks, want) {
		t.Fatalf("uploaded chunks = %#v, want %#v", got.UploadedChunks, want)
	}
}

func TestUploadChunkRejectsMismatchedUploadID(t *testing.T) {
	svc, _ := newTestUploadService(t)
	file := newMultipartFileHeader(t, "chunk", []byte("abc"))

	err := svc.SaveChunk(7, "bad-upload-id", "0123456789abcdef0123456789abcdef", "demo.mp4", 0, file)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != ErrUploadIDMismatch {
		t.Fatalf("error = %v, want %v", err, ErrUploadIDMismatch)
	}
}

func TestUploadMergeFailsWhenAnyChunkMissing(t *testing.T) {
	svc, _ := newTestUploadService(t)
	userID := uint64(7)
	fileMD5 := "0123456789abcdef0123456789abcdef"
	uploadID := GenerateUploadID(fileMD5, userID, "salt")
	if err := svc.SaveChunk(userID, uploadID, fileMD5, "demo.txt", 0, newMultipartFileHeader(t, "0", []byte("hello"))); err != nil {
		t.Fatal(err)
	}

	_, err := svc.Merge(userID, model.UploadMergeRequest{
		UploadID:    uploadID,
		FileMD5:     fileMD5,
		ChunkLength: 2,
	})
	if err == nil {
		t.Fatalf("expected missing chunk error")
	}
}

func TestUploadMergeCreatesFinalFileAndCleansChunks(t *testing.T) {
	svc, dirs := newTestUploadService(t)
	var scheduled []func()
	var scheduledDelay time.Duration
	svc.scheduleDelete = func(delay time.Duration, fn func()) {
		scheduledDelay = delay
		scheduled = append(scheduled, fn)
	}
	userID := uint64(7)
	fileMD5 := "0123456789abcdef0123456789abcdef"
	uploadID := GenerateUploadID(fileMD5, userID, "salt")
	if err := svc.SaveChunk(userID, uploadID, fileMD5, "demo.txt", 0, newMultipartFileHeader(t, "0", []byte("hello "))); err != nil {
		t.Fatal(err)
	}
	if err := svc.SaveChunk(userID, uploadID, fileMD5, "demo.txt", 1, newMultipartFileHeader(t, "1", []byte("world"))); err != nil {
		t.Fatal(err)
	}

	got, err := svc.Merge(userID, model.UploadMergeRequest{
		UploadID:    uploadID,
		FileMD5:     fileMD5,
		ChunkLength: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.URL != "http://test.local/larges/0123456789abcdef0123456789abcdef.txt" || got.Msg != "合并成功" {
		t.Fatalf("unexpected merge response: %#v", got)
	}
	data, err := os.ReadFile(filepath.Join(dirs.large, "0123456789abcdef0123456789abcdef.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Fatalf("final file = %q", string(data))
	}
	if _, err := os.Stat(filepath.Join(dirs.chunks, uploadID)); !os.IsNotExist(err) {
		t.Fatalf("expected chunks dir cleanup, stat err = %v", err)
	}
	if scheduledDelay != 10*time.Minute {
		t.Fatalf("scheduled delay = %v, want %v", scheduledDelay, 10*time.Minute)
	}
	if len(scheduled) != 1 {
		t.Fatalf("scheduled delete count = %d, want 1", len(scheduled))
	}

	repeated, err := svc.Merge(userID, model.UploadMergeRequest{
		UploadID:    uploadID,
		FileMD5:     fileMD5,
		ChunkLength: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if repeated.URL != got.URL || repeated.Msg != got.Msg {
		t.Fatalf("unexpected repeated merge response: %#v", repeated)
	}
	if len(scheduled) != 1 {
		t.Fatalf("repeated merge scheduled delete count = %d, want 1", len(scheduled))
	}

	scheduled[0]()
	if _, err := os.Stat(filepath.Join(dirs.large, "0123456789abcdef0123456789abcdef.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected scheduled delete to remove final file, stat err = %v", err)
	}
}

type testUploadDirs struct {
	large  string
	chunks string
}

func newTestUploadService(t *testing.T) (*UploadService, testUploadDirs) {
	t.Helper()
	tmp := t.TempDir()
	dirs := testUploadDirs{
		large:  filepath.Join(tmp, "larges"),
		chunks: filepath.Join(tmp, "chunks"),
	}
	oldConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		Attach: config.AttachConfig{
			ViewLargeFileBaseURL: "http://test.local/larges",
			UploadLargeFilePath:  dirs.large,
			ChunkDirPath:         dirs.chunks,
			ChunkDirSalt:         "salt",
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = oldConfig
	})
	return NewUploadService(persistence.NewUploadRepository()), dirs
}

func newMultipartFileHeader(t *testing.T, fieldName string, data []byte) *multipart.FileHeader {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fieldName)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(int64(body.Len())); err != nil {
		t.Fatal(err)
	}
	files := req.MultipartForm.File[fieldName]
	if len(files) != 1 {
		t.Fatalf("expected one multipart file, got %d", len(files))
	}
	return files[0]
}
