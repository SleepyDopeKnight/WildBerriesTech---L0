package reader_json

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Open(filesDirectory string) [][]byte {
	files, err := os.ReadDir(filesDirectory)
	if err != nil {
		log.Println(err)
	}

	var filesData [][]byte

	for _, file := range files {
		filePath := filepath.Join(filesDirectory, file.Name())
		if strings.HasSuffix(filePath, ".json") {
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				log.Println(err)
			}

			filesData = append(filesData, fileData)
		}
	}

	return filesData
}
