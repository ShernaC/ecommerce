package tools

import "gorm.io/gorm"

func IsDeletedAtNull(query *gorm.DB) *gorm.DB {
	return query.Where("deleted_at IS NULL")
}
