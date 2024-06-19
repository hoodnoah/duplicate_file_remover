package main

import (
	// std
	"fmt"
	"io"
	"os"

	// internals
	"github.com/hoodnoah/duplicate_file_remover/internal/args"
	"github.com/hoodnoah/duplicate_file_remover/internal/files"

	// community
	"github.com/OneOfOne/xxhash"
)

func main() {
	argsResult, err := args.Consume(os.Args)
	if err != nil {
		fmt.Printf("Failed to consume args with error: %s\n", err)
		fmt.Printf("Exiting.\n")
		panic(1)
	}

	filesIterator, err := files.FilesIteratorNew(argsResult.WorkingPath)
	if err != nil {
		fmt.Printf("Failed to create files iterator with error %s\n", err)
		fmt.Printf("Exiting.\n")
		panic(1)
	}

	var fileHashMaps map[uint64][]string = make(map[uint64][]string, 0)
	hasher := xxhash.New64()
	for filepath := range filesIterator.C {
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

	num_duplicates := 0
	duplicates_removed := 0

	for _, paths := range fileHashMaps {
		if len(paths) > 1 {
			num_duplicates += len(paths) - 1
			parent := paths[0]
			children := paths[1:]

			fmt.Printf("Parent file %s has children:\n", parent)
			for _, child := range children {
				fmt.Printf("\t%s\n", child)

				err = os.Remove(child)
				if err != nil {
					fmt.Printf("Failed to remove file %s", child)
				} else {
					duplicates_removed++
				}
			}
		}
	}

	fmt.Println()
	fmt.Printf("Found %d duplicate files.\n", num_duplicates)
	fmt.Printf("Removed %d duplicate files.\n", duplicates_removed)
}
