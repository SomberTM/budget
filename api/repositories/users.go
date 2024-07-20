package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
	"database/sql"
)

type UsersRepository interface {
	GetUser(ctx context.Context, userId string) (models.User, bool, error)
	GetUserByUserName(ctx context.Context, userName string) (models.User, bool, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type DatabaseUsersRepository struct {
	db dependencies.Database
}

func NewDatabaseUsersRepository(db dependencies.Database) *DatabaseUsersRepository {
	return &DatabaseUsersRepository{db: db}
}

var nilUser models.User = models.User{}

func (r *DatabaseUsersRepository) GetUser(ctx context.Context, userId string) (models.User, bool, error) {
	var user models.User
	row := r.db.GetConnection().QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", userId)
	if err := row.Scan(&user.Id, &user.UserName, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nilUser, false, nil
		}

		return nilUser, false, err
	}

	return user, true, nil
}

func (r *DatabaseUsersRepository) GetUserByUserName(ctx context.Context, userName string) (models.User, bool, error) {
	var user models.User
	row := r.db.GetConnection().QueryRowContext(ctx, "SELECT * FROM users WHERE user_name = $1", userName)
	if err := row.Scan(&user.Id, &user.UserName, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nilUser, false, nil
		}

		return nilUser, false, err
	}

	return user, true, nil
}

func (r *DatabaseUsersRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	row := r.db.GetConnection().QueryRowContext(ctx, "INSERT INTO users (id, user_name, password_hash) VALUES (default, $1, $2) RETURNING id", user.UserName, user.PasswordHash)
	if err := row.Scan(&user.Id); err != nil {
		return nilUser, err
	}

	return user, nil
}

type NilUsersRepository struct{}

func (r *NilUsersRepository) GetUser(ctx context.Context, userId string) (models.User, bool, error) {
	return nilUser, false, nil
}
func (r *NilUsersRepository) GetUserByUserName(ctx context.Context, userName string) (models.User, bool, error) {
	return nilUser, false, nil
}
func (r *NilUsersRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return nilUser, nil
}
