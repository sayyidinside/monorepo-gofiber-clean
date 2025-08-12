package helpers

import (
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/model"
	"gorm.io/gorm"
)

func Filter(query *model.QueryGet, allowedFields map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		filterBy := query.FilterBy
		filterValue := query.Filter

		// Validate the filter_by field and retrieve the corresponding database field
		dbField, isValidFilterField := allowedFields[filterBy]

		if isValidFilterField && filterValue != "" { // when option valid filtered data from db
			db = db.Where(dbField+" = ?", filterValue)
		}

		return db
	}
}
