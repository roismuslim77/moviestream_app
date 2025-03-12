package movie

import (
	"context"
	"net/http"
	"simple-go/application/entity"
	"simple-go/pkg/response"
	"time"
)

type Repository interface {
	GetMovies(ctx context.Context, filter FilterMovie) ([]entity.AllMovie, int64, float64, error)
	GetMovieById(ctx context.Context, id int) (entity.Movie, error)
	CreateMovie(ctx context.Context, req entity.Movie) (entity.Movie, error)
	UpdateMovie(ctx context.Context, req entity.Movie, id int) (entity.Movie, error)
	DeleteMovie(ctx context.Context, id int) error

	CreateMovieView(ctx context.Context, req entity.MovieView) (entity.MovieView, error)
	GetCustomerMovieVoteById(ctx context.Context, id, customerId int) (entity.MovieVote, error)
	CreateMovieVote(ctx context.Context, req entity.MovieVote) (entity.MovieVote, error)
	DeleteMovieVote(ctx context.Context, id, customerId int) error

	CountMovieView(ctx context.Context, movieId int, resp chan<- int64)
	CountMovieVote(ctx context.Context, movieId int, resp chan<- int64)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) service {
	return service{
		repository: repo,
	}
}

func (s service) GetMovies(ctx context.Context, customerId int, filter FilterMovie) ([]MoviesResponse, PaginateListing, response.ErrorResponse) {
	var resp []MoviesResponse
	var paginate PaginateListing

	filter.CustomerId = customerId
	listing, rows, page, err := s.repository.GetMovies(ctx, filter)
	if err != nil {
		errResp := response.Error("22101").WithStatusCode(http.StatusInternalServerError)
		return nil, paginate, *errResp
	}

	for _, list := range listing {
		var movie MoviesResponse
		movie.Title = list.Title
		movie.Description = list.Description
		movie.Duration = list.Duration
		movie.Artist = list.Artist
		movie.Genre = list.Genre
		movie.WatchUrl = list.WatchUrl
		movie.UpdatedAt = list.UpdatedAt
		movie.CreatedAt = list.CreatedAt

		movie.IsViewed = false
		if list.IsViewed != nil {
			movie.IsViewed = true
		}

		movie.IsVoted = false
		if list.IsVoted != nil {
			movie.IsVoted = true
		}

		countVote := make(chan int64)
		countView := make(chan int64)
		go s.repository.CountMovieView(ctx, list.ID, countView)
		go s.repository.CountMovieVote(ctx, list.ID, countVote)
		votes := <-countVote
		views := <-countView

		movie.TotalVote = votes
		movie.TotalViews = views
		resp = append(resp, movie)
	}

	paginate.TotalData = rows
	paginate.TotalPage = page
	return resp, paginate, *response.NotError()
}

func (s service) WatchMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse {
	movie, err := s.repository.GetMovieById(ctx, movieId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if movie.IsEmpty {
		return *response.Error("22101").WithError("movie not found").WithStatusCode(http.StatusNotFound)
	}

	payloadMovie := entity.MovieView{
		MovieId:    movieId,
		CustomerId: customerId,
		CreatedAt:  time.Now(),
	}
	_, err = s.repository.CreateMovieView(ctx, payloadMovie)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}

	return *response.NotError()
}

func (s service) VoteMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse {
	movie, err := s.repository.GetMovieById(ctx, movieId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if movie.IsEmpty {
		return *response.Error("22101").WithError("movie not found").WithStatusCode(http.StatusNotFound)
	}

	customerVote, err := s.repository.GetCustomerMovieVoteById(ctx, movieId, customerId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if !customerVote.IsEmpty {
		return *response.Error("22101").WithError("movie has been voted").WithStatusCode(http.StatusBadRequest)
	}

	payloadMovie := entity.MovieVote{
		MovieId:    movieId,
		CustomerId: customerId,
		CreatedAt:  time.Now(),
	}
	_, err = s.repository.CreateMovieVote(ctx, payloadMovie)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}

	return *response.NotError()
}

func (s service) UnVoteMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse {
	movie, err := s.repository.GetMovieById(ctx, movieId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if movie.IsEmpty {
		return *response.Error("22101").WithError("movie not found").WithStatusCode(http.StatusNotFound)
	}

	customerVote, err := s.repository.GetCustomerMovieVoteById(ctx, movieId, customerId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if customerVote.IsEmpty {
		return *response.Error("22101").WithError("movie has been un voted").WithStatusCode(http.StatusBadRequest)
	}

	err = s.repository.DeleteMovieVote(ctx, movieId, customerId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}

	return *response.NotError()
}

func (s service) AdminCreateMovie(ctx context.Context, customerId int, req CreateMovieRequest) response.ErrorResponse {
	payloadMovie := entity.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Artist:      req.Artist,
		Genre:       req.Genre,
		WatchUrl:    req.WatchUrl,
		CreatedAt:   time.Now(),
		CreatedId:   customerId,
	}
	_, err := s.repository.CreateMovie(ctx, payloadMovie)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}

	return *response.NotError()
}

func (s service) AdminUpdateMovie(ctx context.Context, customerId, movieId int, req UpdateMovieRequest) response.ErrorResponse {
	movie, err := s.repository.GetMovieById(ctx, movieId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if movie.IsEmpty {
		return *response.Error("22101").WithError("movie not found").WithStatusCode(http.StatusNotFound)
	}

	movie.Title = req.Title
	movie.Description = req.Description
	movie.Duration = req.Duration
	movie.Artist = req.Artist
	movie.Genre = req.Genre
	movie.WatchUrl = req.WatchUrl
	movie.UpdatedAt = time.Now()
	movie.CreatedId = customerId
	_, err = s.repository.UpdateMovie(ctx, movie, movieId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	return *response.NotError()
}

func (s service) AdminDeleteMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse {
	movie, err := s.repository.GetMovieById(ctx, movieId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusBadRequest)
	}
	if movie.IsEmpty {
		return *response.Error("22101").WithError("movie not found").WithStatusCode(http.StatusNotFound)
	}

	err = s.repository.DeleteMovie(ctx, movieId)
	if err != nil {
		return *response.Error("22101").WithError(err.Error()).WithStatusCode(http.StatusInternalServerError)
	}

	return *response.NotError()
}
