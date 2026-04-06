package data

import (
	"testing"

	"greenlight.isez.dev/internal/assert"
)

func TestMovieModel_GetAll(t *testing.T) {
	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafeList: []string{"id", "title", "year", "runtime"},
	}

	tests := []struct {
		name      string
		title     string
		wantCount int
		wantTitle string
	}{
		{
			name:      "Found by title",
			title:     "Black",
			wantCount: 1,
			wantTitle: "Black Panther",
		},
		{
			name:      "Not found",
			title:     "Nonexistent",
			wantCount: 0,
		},
		{
			name:      "Empty filter returns all",
			title:     "",
			wantCount: 1,
			wantTitle: "Black Panther",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := MovieModel{DB: db}

			movies, _, err := m.GetAll(tt.title, []string{}, filters)

			assert.NilError(t, err)
			assert.Equal(t, len(movies), tt.wantCount)
			if tt.wantTitle != "" {
				assert.Equal(t, movies[0].Title, tt.wantTitle)
			}
		})
	}
}
