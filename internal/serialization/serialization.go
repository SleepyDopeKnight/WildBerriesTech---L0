package serialization

import (
	"L0/internal/database/models"
	"encoding/json"
	"log"
)

func FileDeserialize(fileData []byte) *models.Orders {
	var orders models.Orders
	err := json.Unmarshal(fileData, &orders)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &orders
}
