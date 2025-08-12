package helpers

import (
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
	"gorm.io/gorm"
)

func Order(query *model.QueryGet, allowedFields map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Ordering logic
		orderBy := query.OrderBy
		orderValue := query.Order

		// Validate the order_by field and retrieve the corresponding database field
		dbField, isValidOrderField := allowedFields[orderBy]

		if isValidOrderField { // when option valid ordered data from db
			if orderValue != "asc" && orderValue != "desc" {
				orderValue = "asc"
			}
			db = db.Order(dbField + " " + orderValue)
		}

		return db
	}
}
