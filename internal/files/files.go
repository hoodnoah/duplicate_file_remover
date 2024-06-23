package files

import (
	"io/fs"
	"path/filepath"
)

type FilesIterator struct {
	path string
	C    chan string
}

// abstraction for testing
type IDirWalker interface {
	WalkDir(string, fs.WalkDirFunc) error
}

type DefaultDirWalker struct{}

func (dw *DefaultDirWalker) WalkDir(path string, pathFunc fs.WalkDirFunc) error {
	return filepath.WalkDir(path, pathFunc)
}

func FilesIteratorNew(path string, numCores int, dirWalker IDirWalker) (*FilesIterator, error) {
	filesChannel := make(chan string, numCores)

	filesIterator := &FilesIterator{
		path: path,
		C:    filesChannel,
	}

	go func() {
		err := populateFilesChannel(filesIterator.path, filesIterator.C, dirWalker)
		if err != nil {
			close(filesIterator.C)
		}
	}()

	return filesIterator, nil
}

// fills the channel with files, not including pure directories
func populateFilesChannel(path string, filesChannel chan string, dirWalker IDirWalker) error {
	defer close(filesChannel)

	err := dirWalker.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			filesChannel <- path
		}

		return nil
	})
	return err
}
