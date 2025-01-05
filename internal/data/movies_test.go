package data

import (
	"testing"

	"greenlight.isez.dev/internal/assert"
)

func TestListMoviesHandler(t *testing.T) {
	tests := []struct {
		name    string
		movieID int
		want    []Movie
	}{
		{
			name:    "Valid ID",
			movieID: 1,
			want: []Movie{{
				Title:   "Black Panther",
				Year:    2018,
				Runtime: 134,
				Genres:  []string{"action", "adventure"},
				Version: 1},
			},
		}, {
			name:    "Zero ID",
			movieID: 0,
			want:    []Movie{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := MovieModel{DB: db}

			movies, _, err := m.GetAll("", []string{""}, Filters{Page: 1, PageSize: 20, Sort: "id", SortSafeList: []string{"id", "title", "year", "runtime"}})

			assert.Equal(t, movies[0].Title, tt.want[0].Title)
			assert.NilError(t, err)
		})
	}
}
