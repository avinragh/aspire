package web

type FindLoansParams struct {
	UserID *int64  `json:"userId,omitempty"`
	State  *string `json:"state,omitempty"`
	Sort   *string `json:"sort,omitempty"`
	Limit  *int64  `json:"limit,omitempty"`
	Page   *int64  `json:"page,omitempty"`
}

type FindInstallmentsParams struct {
	LoanID *int64  `json:"userId,omitempty"`
	State  *string `json:"state,omitempty"`
	Sort   *string `json:"sort,omitempty"`
	Limit  *int64  `json:"limit,omitempty"`
	Page   *int64  `json:"page,omitempty"`
}
