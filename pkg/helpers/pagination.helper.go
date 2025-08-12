package helpers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/interfaces/model"
	"gorm.io/gorm"
)

func Paginate(query *model.QueryGet) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(query.Page)
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(query.Limit)
		if limit <= 0 {
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func GeneratePaginationMetadata(query *model.QueryGet, url string, totalData int64) *Pagination {
	// initilize required variable
	var nextPage, previousPage string
	var fromRow, toRow int
	totalRow := int(totalData)

	// getting and setting page
	page, _ := strconv.Atoi(query.Page)
	if page <= 0 {
		page = 1
	}

	// getting and setting page
	limit, _ := strconv.Atoi(query.Limit)
	if limit <= 0 {
		limit = 10
	}

	// Calculate total page using totalRow [len(data)] and limit
	totalPages := int(math.Ceil(float64(totalRow) / float64(limit)))

	// Set url for first and last page
	// firstPage := fmt.Sprintf("%s?page=1&limit=%d", *url, limit)
	// lastPage := fmt.Sprintf("%s?page=%d&limit=%d", *url, totalPages, limit)

	// Set url for current, previous, and next page
	currentPage := fmt.Sprintf("%s?page=%d&limit=%d", url, page, limit)
	if page > 1 {
		previousPage = fmt.Sprintf("%s?page=%d&limit=%d", url, page-1, limit)
	}
	if page < totalPages {
		nextPage = fmt.Sprintf("%s?page=%d&limit=%d", url, page+1, limit)
	}

	// Set from and to row (index)
	if page == 1 {
		fromRow = 1
		if limit > totalRow {
			toRow = totalRow
		} else {
			toRow = limit
		}
	} else if page <= totalPages {
		fromRow = ((page - 1) * limit) + 1

		if page == totalPages {
			toRow = totalRow
		} else {
			toRow = page * limit
		}
	}

	return &Pagination{
		CurrentPage: page,
		TotalItems:  totalRow,
		TotalPages:  totalPages,
		ItemPerPage: limit,
		Self:        currentPage,
		Prev:        &previousPage,
		Next:        &nextPage,
		FromRow:     fromRow,
		ToRow:       toRow,
	}
}
