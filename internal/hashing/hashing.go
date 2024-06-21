package hashing

import (
	"bufio"
	"fmt"
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
	var waitGroup sync.WaitGroup

	for range h.numThreads {
		waitGroup.Add(1)
		go func() {
			hasher := xxhash.New64()
			for filePath := range h.filePathsChannel {
				hashedFile, err := hashFile(hasher, filePath)
				if err != nil {
					fmt.Printf("Failed to hash file %s with error %s, skipping.\n", filePath, err)
				} else {
					h.hashMapMutex.Lock()
					h.fileHashMap[hashedFile.hash] = append(h.fileHashMap[hashedFile.hash], hashedFile.path)
					h.hashMapMutex.Unlock()
				}
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait() // let execution run
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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// use buffered reader to reduce system calls per flame graph
	bufferedReader := bufio.NewReaderSize(file, 16*1024) // 16KiB buffer size
	if _, err := bufferedReader.WriteTo(hasher); err != nil {
		return nil, err
	}

	hash := hasher.Sum64()
	defer hasher.Reset()

	return &hashedFile{
		hash: hash,
		path: path,
	}, nil
}
