package args

import (
	"path/filepath"
)

type Args struct {
	WorkingPath string
}

type AbsPathResolver interface { // abstract interface, for testing
	Abs(path string) (string, error)
}

type DefaultPathResolver struct{}

func (dr *DefaultPathResolver) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

// consumes the args, putting them into a struct
func Consume(args []string, resolver AbsPathResolver) (*Args, error) {
	var pathString string

	if len(args) < 2 { // default to current directory
		pathString = "."
	} else {
		pathString = args[1]
	}

	path, err := resolver.Abs(pathString)

	if err != nil {
		return nil, err
	}

	return &Args{
		WorkingPath: path,
	}, nil
}
