package qonto

// OrganizationResponse holds the data returned by the Qonto API when fetching organization info.
type OrganizationResponse struct {
	Organization *Organization `json:"organization"`
}

// TransactionsResponse holds the data returned by the Qonto API when fetching a page of transactions.
type TransactionsResponse struct {
	Transactions []*Transaction `json:"transactions"`
	Meta         *ResponseMeta  `json:"meta"`
}

// ResponseMeta holds the paging information when fetching a page of transactions.
type ResponseMeta struct {
	CurrentPage int  `json:"current_page"`
	NextPage    *int `json:"next_page"`
	PrevPage    *int `json:"prev_page"`
	TotalPages  int  `json:"total_pages"`
	TotalCount  int  `json:"total_count"`
	PerPage     int  `json:"per_page"`
}

// IsLastPage returns true when the meta-data corresponds to the last page.
func (m *ResponseMeta) IsLastPage() bool {
	return m.NextPage == nil
}

// IsLastPage returns true when the meta-data corresponds to the first page.
func (m *ResponseMeta) IsFirstPage() bool {
	return m.PrevPage == nil
}
