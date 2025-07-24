package web

type FilterProduct struct {
	Search     string `query:"search"`
	CategoryID string `query:"category_id"`
	Sort       string `query:"sort"`
	Page       string `query:"page"`
	Limit      string `query:"limit"`
}
