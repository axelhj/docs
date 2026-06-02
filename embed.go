package service

import (
	"embed"
	"fmt"
)

//go:embed schemas/*.json
var SchemasFS embed.FS

//go:embed queries/*.json
var QueriesFS embed.FS

//go:embed test_doc/*.json
var TestDocFS embed.FS

func DoReadFile(loc embed.FS, filename string) ([]byte, error) {
	data, err := loc.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}
	return data, nil
}
