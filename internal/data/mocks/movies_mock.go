package mocks

import (
	"database/sql"
	"time"

	"greenlight.isez.dev/internal/data"
)

type MovieModel_Mock struct {
	DB *sql.DB
}

func (m MovieModel_Mock) Insert(movie *data.Movie) error {

	return nil
}

func (m MovieModel_Mock) Get(id int64) (*data.Movie, error) {
	test_movie := data.Movie{
		ID:        1,
		CreatedAt: time.Now(),
		Title:     "Moana",
		Year:      2024,
		Runtime:   120,
		Genres:    []string{"Action"},
		Version:   1,
	}

	return &test_movie, nil
}

func (m MovieModel_Mock) Update(movie *data.Movie) error {

	return nil
}

func (m MovieModel_Mock) Delete(id int64) error {

	return nil
}

func (m MovieModel_Mock) GetAll(title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
	test_movies := []*data.Movie{{
		ID:        1,
		CreatedAt: time.Now(),
		Title:     "Moana",
		Year:      2024,
		Runtime:   120,
		Genres:    []string{"Action"},
		Version:   1}, {
		ID:        2,
		CreatedAt: time.Now(),
		Title:     "Moana",
		Year:      2024,
		Runtime:   120,
		Genres:    []string{"Action"},
		Version:   1,
	},
	}

	metadata := data.Metadata{
		CurrentPage:  1,
		PageSize:     20,
		FirstPage:    1,
		LastPage:     1,
		TotalRecords: 2,
	}

	return test_movies, metadata, nil
}
