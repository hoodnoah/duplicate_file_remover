package args_test

import (
	"strings"
	"testing"

	"github.com/hoodnoah/duplicate_file_remover/internal/args"
)

func TestArgs(t *testing.T) {
	t.Run("gets the absolute path correctly", func(t *testing.T) {
		testArguments := []string{"dfr", "/path/to/files"}

		argsResult, err := args.Consume(testArguments)
		if err != nil {
			t.Error(err)
		}

		if strings.Compare(argsResult.WorkingPath, testArguments[1]) != 0 {
			t.Errorf("expected argsResult.WorkingPath to equal %s, got %s instead", testArguments[1], argsResult.WorkingPath)
		}
	})
}
