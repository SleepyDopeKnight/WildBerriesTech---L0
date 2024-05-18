package serialization

import (
	"encoding/json"
	"log"

	"L0/internal/database/models"
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
