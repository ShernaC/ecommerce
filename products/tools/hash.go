package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"products/model"
)

func GenerateSKU(product *model.Product) string {
	data := fmt.Sprintf("%d-%d-%s-%d",
		product.ID,
		product.SellerID,
		product.Name,
		product.CreatedAt.Unix(),
	)

	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])

	return fmt.Sprintf("SKU-%s", hashStr[:12])
}
