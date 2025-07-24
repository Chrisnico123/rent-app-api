package domain

import (
	"fmt"
	"rent-application/shared/helper"
	"rent-application/shared/web"
	"strconv"
)

type Banner struct {
	Id   string
	Name string
	Img  string
}

type FilterBannerPagination struct {
	Search     string
	Pagination Pagination
}

func ToDomainFilterBannerPagination(q web.FilterBannerPagination) FilterBannerPagination {
	pageInt, _ := strconv.Atoi(q.Page)
	limitInt, _ := strconv.Atoi(q.Limit)

	var page, limit string

	if q.Page != "" {
		if q.Limit != "" {
			page = fmt.Sprintf("%d", (pageInt-1)*limitInt)
			limit = fmt.Sprintf("%d", limitInt)
		} else {
			page = fmt.Sprintf("%d", (pageInt-1)*5)
			limit = fmt.Sprintf("%d", 5)
		}
	} else {
		if q.Limit != "" {
			page = "1"
			limit = fmt.Sprintf("%d", limitInt)
		}
	}

	return FilterBannerPagination{
		Search: q.Search,
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
		},
	}
}

func (q *FilterBannerPagination) QueryBuildFilterSearchPagination() (filter, pagination string) {
	filter = BuildBannerQuerySearch(q.Search)

	if q.Pagination.Page != "" && q.Pagination.Limit != "" {
		pagination = fmt.Sprintf("OFFSET %s LIMIT %s", q.Pagination.Page, q.Pagination.Limit)
	}
	return
}

func BuildBannerQuerySearch(search string) (filter string) {
	fieldsAndValues := map[string]string{
		"b.name": search,
		// Add more field-value pairs as needed
	}

	filter = helper.BuildManyColumnFilter(fieldsAndValues, filter)

	return
}
