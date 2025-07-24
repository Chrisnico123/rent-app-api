package web

type FilterBannerPagination struct {
	Search string `query:"search"`
	Limit  string `query:"limit"`
	Page   string `query:"page"`
}

type Banner struct {
	Id   string
	Name string
	Img  string
}
