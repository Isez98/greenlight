package mocks

import (
	"database/sql"

	"greenlight.isez.dev/internal/data"
)

type PermissionModel_Mock struct {
	DB *sql.DB
}

func (m PermissionModel_Mock) GetAllForUser(userID int64) (data.Permissions, error) {
	return data.Permissions{"movies:read"}, nil
}

func (m PermissionModel_Mock) AddForUser(userID int64, codes ...string) error {

	return nil
}
