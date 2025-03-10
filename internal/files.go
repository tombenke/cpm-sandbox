package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetFilePathWithoutExtension(filePath string) string {
	// Get the path to the file
	path := filepath.Dir(filePath)

	// Get filename without path
	fileNameWithExt := filepath.Base(filePath)

	// Remove extension
	fileName := strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))

	return filepath.Join(path, fileName)
}

// LoadFile loads []byte content from a file
func LoadFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return content
}

// WriteFile creates a file and writes the content into it
func WriteFile(fileName string, content string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// A file egy io.Writer interfészt implementáló típus
	writer := file

	fmt.Fprintln(writer, content)
}
