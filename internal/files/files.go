package files

import (
	"io/fs"
	"path/filepath"
)

type FilesIterator struct {
	path string
	C    chan string
}

func FilesIteratorNew(path string, numCores int) (*FilesIterator, error) {
	filesChannel := make(chan string, numCores)

	filesIterator := &FilesIterator{
		path: path,
		C:    filesChannel,
	}

	go func() {
		err := populateFilesChannel(filesIterator.path, filesIterator.C)
		if err != nil {
			close(filesIterator.C)
		}
	}()

	return filesIterator, nil
}

// fills the channel with files, not including pure directories
func populateFilesChannel(path string, filesChannel chan string) error {
	defer close(filesChannel)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
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
