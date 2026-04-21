package services

type PocketBaseErrorResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

type PocketBaseListResponse struct {
	Page       int              `json:"page"`
	PerPage    int              `json:"perPage"`
	TotalItems int              `json:"totalItems"`
	TotalPages int              `json:"totalPages"`
	Items      []map[string]any `json:"items"`
}

type PocketBaseListOptions struct {
	Page      int
	PerPage   int
	Sort      string
	Filter    string
	Expand    string
	Fields    string
	SkipTotal bool
}
