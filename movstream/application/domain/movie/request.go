package movie

type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Duration    int64  `json:"duration" binding:"required"`
	Artist      string `json:"artist" binding:"required"`
	Genre       string `json:"genre" binding:"required"`
	WatchUrl    string `json:"watch_url" binding:"required"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int64  `json:"duration"`
	Artist      string `json:"artist"`
	Genre       string `json:"genre"`
	WatchUrl    string `json:"watch_url"`
}

type FilterMovie struct {
	Page       string `json:"page"`
	Limit      string `json:"limit"`
	Search     string `json:"search"`
	CustomerId int    `json:"customer_id"`
}
