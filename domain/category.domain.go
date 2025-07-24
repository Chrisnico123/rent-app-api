package domain

import (
	"fmt"
	"rent-application/shared/helper"
	"rent-application/shared/web"
	"strconv"
)

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

type FilterSearchPagination struct {
	Search     string
	Pagination Pagination
}

func ToDomainFilterSearchPagination(q web.FilterSearchPagination) FilterSearchPagination {
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

	return FilterSearchPagination{
		Search: q.Search,
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
		},
	}
}

func (q *FilterSearchPagination) QueryBuildFilterSearchPagination() (filter, pagination string) {
	filter = BuildPresenceQuerySearch(q.Search)

	if q.Pagination.Page != "" && q.Pagination.Limit != "" {
		pagination = fmt.Sprintf("OFFSET %s LIMIT %s", q.Pagination.Page, q.Pagination.Limit)
	}
	return
}

func BuildPresenceQuerySearch(search string) (filter string) {
	fieldsAndValues := map[string]string{
		"c.name": search,
		// Add more field-value pairs as needed
	}

	filter = helper.BuildManyColumnFilter(fieldsAndValues, filter)

	return
}
