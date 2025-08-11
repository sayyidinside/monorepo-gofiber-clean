package model

import "github.com/microcosm-cc/bluemonday"

type (
	QueryGet struct {
		Page     string `query:"page"`
		Limit    string `query:"limit"`
		OrderBy  string `query:"order_by"`
		Order    string `query:"order"`
		FilterBy string `query:"filter_by"`
		Filter   string `query:"filter"`
		SearchBy string `query:"search_by"`
		Search   string `query:"search"`
	}
)

func SanitizeQueryGet(query *QueryGet) {
	sanitizer := bluemonday.StrictPolicy()
	query.Page = sanitizer.Sanitize(query.Page)
	query.Limit = sanitizer.Sanitize(query.Limit)
	query.OrderBy = sanitizer.Sanitize(query.OrderBy)
	query.Order = sanitizer.Sanitize(query.Order)
	query.FilterBy = sanitizer.Sanitize(query.FilterBy)
	query.Filter = sanitizer.Sanitize(query.Filter)
	query.SearchBy = sanitizer.Sanitize(query.SearchBy)
	query.Search = sanitizer.Sanitize(query.Search)
}
