package example

type ListResponse struct {
	Items  []ItemResponse `json:"items"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Total  int            `json:"total"`
}

type ItemResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
