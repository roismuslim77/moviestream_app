package movie

import "time"

type PaginateListing struct {
	TotalData int64   `json:"total_data"`
	TotalPage float64 `json:"total_page"`
}

type MoviesResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int64     `json:"duration"`
	Artist      string    `json:"artist"`
	Genre       string    `json:"genre"`
	WatchUrl    string    `json:"watch_url"`
	TotalViews  int64     `json:"total_views"`
	TotalVote   int64     `json:"total_vote"`
	IsViewed    bool      `json:"is_viewed"`
	IsVoted     bool      `json:"is_voted"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}
