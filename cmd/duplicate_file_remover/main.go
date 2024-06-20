package main

import (
	// std
	"fmt"
	"os"
	"runtime"

	// internals
	"github.com/hoodnoah/duplicate_file_remover/internal/args"
	"github.com/hoodnoah/duplicate_file_remover/internal/files"
	"github.com/hoodnoah/duplicate_file_remover/internal/hashing"
	// community
)

func main() {
	numThreads := runtime.NumCPU()

	argsResult, err := args.Consume(os.Args)
	if err != nil {
		fmt.Printf("Failed to consume args with error: %s\n", err)
		fmt.Printf("Exiting.\n")
		panic(1)
	}

	filesIterator, err := files.FilesIteratorNew(argsResult.WorkingPath, numThreads)
	if err != nil {
		fmt.Printf("Failed to create files iterator with error %s\n", err)
		fmt.Printf("Exiting.\n")
		panic(1)
	}

	fileHasher := hashing.HasherNew(numThreads, filesIterator.C)
	fileHasher.HashFiles()
	duplicates := fileHasher.GetDuplicates()

	num_duplicates := 0
	duplicates_removed := 0

	for _, duplicate := range duplicates {
		num_duplicates += len(duplicate.Children)
		for _, child := range duplicate.Children {
			err := os.Remove(child)
			if err != nil {
				fmt.Printf("Failed to remove file %s with error %s, skipping", child, err)
			}
			duplicates_removed++
		}
	}

	fmt.Println()
	fmt.Printf("Found %d duplicate files.\n", num_duplicates)
	fmt.Printf("Removed %d duplicate files.\n", duplicates_removed)
}
