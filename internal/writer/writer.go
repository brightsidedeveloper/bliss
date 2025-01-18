package writer

import (
	"os"
	"path/filepath"
)

func WriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}
