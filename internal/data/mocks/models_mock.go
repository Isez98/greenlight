package mocks

import (
	"greenlight.isez.dev/internal/data"
)

func TestModels_Mock() data.Models {
	return data.Models{
		Movies:      MovieModel_Mock{},
		Permissions: PermissionModel_Mock{},
		Tokens:      TokenModel_Mock{},
		Users:       UserModel_Mock{},
	}
}
