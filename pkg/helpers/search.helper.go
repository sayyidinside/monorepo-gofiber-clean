package helpers

import (
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/model"
	"gorm.io/gorm"
)

func Search(query *model.QueryGet, allowedFields map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchBy := query.SearchBy
		searchValue := "%" + query.Search + "%"

		// Validate the search_by field and retrieve the corresponding database field
		dbField, isValidSearchField := allowedFields[searchBy]

		if isValidSearchField && searchValue != "" { // when option valid search data from db
			db = db.Where(dbField+" LIKE ?", searchValue)
		}

		return db
	}
}
