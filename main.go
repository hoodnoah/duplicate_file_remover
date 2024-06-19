package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/OneOfOne/xxhash"
)

func main() {
	args := os.Args[1:]

	relativePath := args[0]

	directory, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Printf("Failed to locate path of files to de-duplicate: %s", err)
		panic(1)
	}

	files := make([]string, 0)

	err = filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failed to visit %s...\n", path)
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Failed to traverse files, %s\n", err)
	}

	var fileHashMaps map[uint64][]string = make(map[uint64][]string, 0)

	hasher := xxhash.New64()
	for _, filepath := range files {
		// open the file
		reader, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Failed to open %s with error %s, skipping", filepath, err)
			continue
		}
		io.Copy(hasher, reader)
		hash := hasher.Sum64()
		children := append(fileHashMaps[hash], filepath)
		fileHashMaps[hash] = children
		reader.Close()
		hasher.Reset()
	}

	for _, paths := range fileHashMaps {
		if len(paths) > 1 {
			parent := paths[0]
			children := paths[1:]

			fmt.Printf("Parent file %s has children:\n", parent)
			for _, child := range children {
				fmt.Printf("\t%s\n", child)
			}
		}
	}
}
