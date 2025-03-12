package movie

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"simple-go/application/entity"
	"strconv"
	"strings"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) GetMovies(ctx context.Context, filter FilterMovie) ([]entity.AllMovie, int64, float64, error) {
	var data []entity.AllMovie
	var totalData int64
	var totalPage float64

	//mapping filter data
	where := "1=1"
	if filter.Search != "" {
		search := strings.ToLower(filter.Search)
		where += " and (LOWER(mv.title) like '%" + search + "%' or LOWER(mv.description) like '%" + search + "%' or LOWER(mv.artist) = '" + search + "' or LOWER(mv.genre) = '" + search + "')"
	}

	//pagination
	offset := 0
	size := 10
	if filter.Limit != "" {
		size, _ = strconv.Atoi(filter.Limit)
	}
	if filter.Page != "" {
		page, _ := strconv.Atoi(filter.Page)
		offset = (page - 1) * size
	}

	r.db.
		Table("movies mv").
		Select("mv.id").
		Joins(fmt.Sprintf("LEFT JOIN LATERAL (SELECT b.id FROM movie_views b WHERE b.movie_id = mv.id and b.customer_id = %v ORDER BY b.created_at DESC LIMIT 1) views ON TRUE", filter.CustomerId)).
		Joins(fmt.Sprintf("LEFT JOIN movie_votes votes on votes.movie_id = mv.id and votes.customer_id = %v", filter.CustomerId)).
		Where(where).
		Count(&totalData)

	result := r.db.
		Table("movies mv").
		Select("mv.id, mv.title, mv.description, mv.duration, mv.artist, mv.genre, mv.watch_url, mv.created_at, mv.updated_at, views.id as is_viewed, votes.id as is_voted").
		Joins(fmt.Sprintf("LEFT JOIN LATERAL (SELECT b.id FROM movie_views b WHERE b.movie_id = mv.id and b.customer_id = %v ORDER BY b.created_at DESC LIMIT 1) views ON TRUE", filter.CustomerId)).
		Joins(fmt.Sprintf("LEFT JOIN movie_votes votes on votes.movie_id = mv.id and votes.customer_id = %v", filter.CustomerId)).
		Where(where).
		Offset(offset).Limit(size).
		Order("mv.created_at DESC").
		Find(&data)

	if totalData > 0 {
		totalPage = float64(totalData) / float64(size)
		totalPage = math.Ceil(totalPage)
	} else {
		totalPage = 1
	}

	if result.Error != nil {
		return nil, 0, 1, result.Error
	}
	return data, totalData, totalPage, nil
}

func (r repository) CountMovieVote(ctx context.Context, movieId int, resp chan<- int64) {
	var count int64

	where := fmt.Sprintf("votes.movie_id = %v", movieId)
	r.db.
		Table("movie_votes votes").
		Select("votes.id").
		Where(where).
		Count(&count)

	resp <- count
}

func (r repository) CountMovieView(ctx context.Context, movieId int, resp chan<- int64) {
	var count int64

	where := fmt.Sprintf("views.movie_id = %v", movieId)
	r.db.
		Table("movie_views views").
		Select("views.id").
		Where(where).
		Count(&count)

	resp <- count
}

func (r repository) GetMovieById(ctx context.Context, id int) (entity.Movie, error) {
	var data entity.Movie
	result := r.db.Where("id = ?", id).First(&data)
	if result.Error != nil {
		if result.RowsAffected < 1 {
			data.IsEmpty = true
			return data, nil
		}
		return data, result.Error
	}

	return data, nil
}

func (r repository) CreateMovie(ctx context.Context, req entity.Movie) (entity.Movie, error) {
	if err := r.db.Clauses(clause.Returning{}).Create(&req).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = errors.New("movie is already exist")
			return req, err
		}
		return req, err
	}

	return req, nil
}

func (r repository) UpdateMovie(ctx context.Context, req entity.Movie, id int) (entity.Movie, error) {
	result := r.db.Clauses(&clause.Returning{}).Where("id = ?", id).Updates(&req)
	if result.Error != nil {
		if result.Error == gorm.ErrDuplicatedKey {
			result.Error = errors.New("movie is already exist")
			return req, result.Error
		}
		return req, result.Error
	}

	if result.RowsAffected < 1 {
		return req, errors.New("failed to update movie")
	}

	return req, nil
}

func (r repository) DeleteMovie(ctx context.Context, id int) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Movie{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected < 1 {
		return errors.New("failed to delete movie")
	}

	return nil
}

func (r repository) CreateMovieView(ctx context.Context, req entity.MovieView) (entity.MovieView, error) {
	if err := r.db.Clauses(clause.Returning{}).Create(&req).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = errors.New("movie is already exist")
			return req, err
		}
		return req, err
	}

	return req, nil
}

func (r repository) GetCustomerMovieVoteById(ctx context.Context, id, customerId int) (entity.MovieVote, error) {
	var data entity.MovieVote
	result := r.db.Where("movie_id = ? and customer_id = ?", id, customerId).First(&data)
	if result.Error != nil {
		if result.RowsAffected < 1 {
			data.IsEmpty = true
			return data, nil
		}
		return data, result.Error
	}

	return data, nil
}

func (r repository) CreateMovieVote(ctx context.Context, req entity.MovieVote) (entity.MovieVote, error) {
	if err := r.db.Clauses(clause.Returning{}).Create(&req).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = errors.New("movie is already exist")
			return req, err
		}
		return req, err
	}

	return req, nil
}

func (r repository) DeleteMovieVote(ctx context.Context, id, customerId int) error {
	result := r.db.Where("movie_id = ? and customer_id = ?", id, customerId).Delete(&entity.MovieVote{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected < 1 {
		return errors.New("failed to delete movie")
	}

	return nil
}
