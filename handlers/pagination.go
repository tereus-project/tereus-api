package handlers

type PaginatedResponse[T any] struct {
	Items []T           `json:"items"`
	Meta  PaginatedMeta `json:"meta"`
}

type PaginatedMeta struct {
	ItemCount    int `json:"item_count"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
}
