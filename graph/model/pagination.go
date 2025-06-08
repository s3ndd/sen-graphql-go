package model

// Pagination structure that includes json representation
type Pagination struct {
	PageIndex        int   `json:"page_index"`
	PageSize         int   `json:"page_size"`
	TotalRecordCount int64 `json:"total_record_count"`
}
