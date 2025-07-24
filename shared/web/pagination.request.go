package web

type FilterSearchPagination struct {
	Search string `query:"search"`
	Page   string `query:"page"`
	Limit  string `query:"limit"`
}
