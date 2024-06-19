package args

import (
	"fmt"
	"path/filepath"
)

type Args struct {
	WorkingPath string
}

// consumes the args, putting them into a struct
func Consume(args []string) (*Args, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("expected 2 args, received %d", len(args))
	}

	path, err := filepath.Abs(args[1])
	if err != nil {
		return nil, err
	}

	return &Args{
		WorkingPath: path,
	}, nil
}
