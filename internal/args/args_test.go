package args_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hoodnoah/duplicate_file_remover/internal/args"
)

type MockPathAbsResolver struct{}
type MockPathAbsResolverFails struct{}

const pathPrefix = "/home/userName/homeDirectory/files"

func (mr *MockPathAbsResolver) Abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	} else {
		return filepath.Join(pathPrefix, path), nil
	}
}

func (mr *MockPathAbsResolverFails) Abs(path string) (string, error) {
	return "", fmt.Errorf("Failed to resolve absolute path")
}

func TestArgs(t *testing.T) {
	mockResolver := MockPathAbsResolver{}
	mockResolverFails := MockPathAbsResolverFails{}

	t.Run("gets the absolute path correctly", func(t *testing.T) {
		testArguments := []string{"ddp", "/path/to/files"}

		argsResult, err := args.Consume(testArguments, &mockResolver)
		if err != nil {
			t.Error(err)
		}

		if strings.Compare(argsResult.WorkingPath, testArguments[1]) != 0 {
			t.Errorf("expected argsResult.WorkingPath to equal %s, got %s instead", testArguments[1], argsResult.WorkingPath)
		}
	})

	t.Run("Defaults to the current working directory when no path is passed", func(t *testing.T) {
		// arrange
		testArguments := []string{"ddp"}
		expectedResult := filepath.Join(pathPrefix, ".")

		argsResult, err := args.Consume(testArguments, &mockResolver)
		if err != nil {
			t.Error(err)
		}

		if strings.Compare(argsResult.WorkingPath, expectedResult) != 0 {
			t.Errorf("expected argsResult.WorkingPath to equal '.', got %s instead", argsResult.WorkingPath)
		}
	})

	t.Run("Propagates filepath error when filepath.Abs fails", func(t *testing.T) {
		// arrange
		testArguments := []string{"ddp"}

		// act
		argsResult, err := args.Consume(testArguments, &mockResolverFails)

		// assert
		if err == nil {
			t.Fatalf("Expected an error to be returned, received nil.")
		}

		if argsResult != nil {
			t.Fatalf("Expected a nil result, received %v", argsResult)
		}
	})
}
