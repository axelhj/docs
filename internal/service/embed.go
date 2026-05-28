package service

import (
	"fmt"
	"embed"
)

//go:embed schemas/*.json
var schemaFS embed.FS

//go:embed queries/*.json
var queryFS embed.FS

func DoReadFile(loc embed.FS, filename string) ([]byte, error) {
	data, err := loc.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed")
	}
	return data, nil
}
