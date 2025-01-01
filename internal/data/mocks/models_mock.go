package mocks

import (
	"database/sql"

	"greenlight.isez.dev/internal/data"
)

func TestModels_Mock(db *sql.DB) data.Models {
	return data.Models{
		Movies:      MovieModel_Mock{DB: db},
		Permissions: PermissionModel_Mock{DB: db},
		Tokens:      TokenModel_Mock{DB: db},
		Users:       UserModel_Mock{DB: db},
	}
}
