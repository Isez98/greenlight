package mocks

import (
	"database/sql"
	"time"

	"greenlight.isez.dev/internal/data"
)

type UserModel_Mock struct {
	DB *sql.DB
}

func (m *UserModel_Mock) GetByEmail(email string) (data.User, error) {
	mockUser := data.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		Activated: true,
		Version:   1,
		Password:  data.AnonymousUser.Password,
	}

	// match, err := mockUser.Password.Matches("pa$$word")
	// if err != nil {
	// 	return data.User{}, err
	// }

	return mockUser, nil
}

func (m *UserModel_Mock) GetForToken(tokenscope string, TokenPlaintext string) (data.User, error) {
	mockUser := data.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		Activated: true,
		Version:   1,
		Password:  data.AnonymousUser.Password,
	}
	return mockUser, nil
}
func (m *UserModel_Mock) Insert(user data.User) error {
	switch user.Email {
	case "dupe@example.com":
		return data.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel_Mock) Update(user data.User) error {
	return nil
}

// func (m *UserModel) Authenticate(email, password string) (int, error) {
// 	if email == "alice@example.com" && password == "pa$$word" {
// 		return 1, nil
// 	}

// 	return 0, data.ErrInvalidCredentials
// }

// func (m *UserModel) Exists(id int) (bool, error) {
// 	switch id {
// 	case 1:
// 		return true, nil
// 	default:
// 		return false, nil
// 	}
// }
