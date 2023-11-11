package helpers

import (
	"io"
	"os"
	"testing"
)

func ReadFile(t *testing.T, filePath string) []byte {
	t.Helper()

	file, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	return data
}
