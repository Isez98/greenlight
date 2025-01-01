package mocks

import (
	"database/sql"

	"greenlight.isez.dev/internal/data"
)

func TestModels_Mock(db *sql.DB) data.Models {
	return data.Models{
		// Movies:      MovieModel{DB: db},
		// Permissions: PermissionModel{DB: db},
		// Tokens:      TokenModel{DB: db},
		// Users:       UserModel_Mock{DB: db},
	}
}
