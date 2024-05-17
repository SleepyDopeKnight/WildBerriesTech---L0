package serialization

import (
	"L0/internal/database/models"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func OpenOrdersJSON(filesDirectory string) [][]byte {
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

func FileDeserialize(fileData []byte) *models.Orders {
	var orders models.Orders
	err := json.Unmarshal(fileData, &orders)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &orders
}
