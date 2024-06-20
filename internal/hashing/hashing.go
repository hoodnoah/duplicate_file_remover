package hashing

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/OneOfOne/xxhash"
)

type Hasher struct {
	numThreads       int
	filePathsChannel chan string
	fileHashMap      map[uint64][]string
	hashMapMutex     sync.Mutex
}

type hashedFile struct {
	hash uint64
	path string
}

type FileDuplicate struct {
	Parent   string
	Children []string
}

func HasherNew(numThreads int, filePathsChannel chan string) *Hasher {
	filesMap := make(map[uint64][]string, 0)

	return &Hasher{
		numThreads:       numThreads,
		filePathsChannel: filePathsChannel,
		fileHashMap:      filesMap,
		hashMapMutex:     sync.Mutex{},
	}
}

func (h *Hasher) HashFiles() {
	xHasher := xxhash.New64()
	for filePath := range h.filePathsChannel {
		hashedFile, err := hashFile(xHasher, filePath)
		if err != nil {
			fmt.Printf("Failed to hash file %s with error %s, skipping\n", filePath, err)
		} else {
			h.hashMapMutex.Lock()
			h.fileHashMap[hashedFile.hash] = append(h.fileHashMap[hashedFile.hash], hashedFile.path)
			h.hashMapMutex.Unlock()
		}
	}
}

func (h *Hasher) GetDuplicates() []FileDuplicate {
	duplicates := make([]FileDuplicate, 0)

	for _, paths := range h.fileHashMap {
		if len(paths) > 1 {
			duplicate := FileDuplicate{
				Parent:   paths[0],
				Children: paths[1:],
			}
			duplicates = append(duplicates, duplicate)
		}
	}

	return duplicates
}

// hashes a file, returning the hash and the path
func hashFile(hasher *xxhash.XXHash64, path string) (*hashedFile, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	io.Copy(hasher, reader)
	hash := hasher.Sum64()
	defer hasher.Reset()

	return &hashedFile{
		hash: hash,
		path: path,
	}, nil
}
