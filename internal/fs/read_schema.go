package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// This works during go run, but fails when the binary is executed from another directory.
func LoadSchema() ([]byte, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed")
	}
	dir := filepath.Dir(filename)
	return os.ReadFile(filepath.Join(dir, "../../schemas/initial.json"))
}

// This works as long as files are located next to the executable. , but fails when the binary is executed from another directory.
func LoadSchemaCwd() ([]byte, error) {
	data, err := os.ReadFile("schemas/initial.json")
	if err != nil {
		return nil, fmt.Errorf("failed")
	}
	return data, nil
}
