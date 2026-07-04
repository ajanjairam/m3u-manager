package models

type Pagination[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
}

func NewPagination[T any](items []T, total int64, page int64, pageSize int64) Pagination[T] {
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}
	return Pagination[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
