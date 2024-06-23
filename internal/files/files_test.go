package files_test

import (
	// test libs
	"io/fs"
	"testing"
)

type mockDirEntry struct {
	name  string
	isDir bool
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

type MockWalkDir struct {
	pathsList []mockDirEntry
}

func (md *MockWalkDir) WalkDir(path string, walkdirFunc fs.WalkDirFunc) error {
	for _, file := range md.pathsList {
		walkdirFunc(file.name, file, nil)
	}

	return nil
}

func TestFiles(t *testing.T) {
	mockPaths := []mockDirEntry{
		mockDirEntry{
			name:  "/path/to/file1.ext",
			isDir: false,
		},
		mockDirEntry{
			name:  "/path/to/file2.ext",
			isDir: false,
		},
		mockDirEntry{
			name:  "/path/to/dir1",
			isDir: true,
		},
		mockDirEntry{
			name:  "path/to/dir1/file1.ext",
			isDir: false,
		},
	}
}
