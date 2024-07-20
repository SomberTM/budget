package repositories

import (
	"budget/api/dependencies"
	"budget/api/models"
	"context"
	"database/sql"
)

type SessionsRepository interface {
	GetSession(ctx context.Context, sessionId string) (models.Session, bool, error)
	GetSessionWithUser(ctx context.Context, sessionId string) (models.Session, models.User, bool, error)
	CreateUserSession(ctx context.Context, userId string) (models.Session, []byte, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type DatabaseSessionsRepository struct {
	db dependencies.Database
}

func NewDatabaseSessionsRepository(db dependencies.Database) *DatabaseSessionsRepository {
	return &DatabaseSessionsRepository{db: db}
}

var nilSession models.Session = models.Session{}

func (r *DatabaseSessionsRepository) GetSession(ctx context.Context, sessionId string) (models.Session, bool, error) {
	var session models.Session
	row := r.db.GetConnection().QueryRowContext(ctx, "SELECT * FROM sessions WHERE id = $1", sessionId)
	if err := row.Scan(&session.Id, &session.UserId, &session.TokenHash); err != nil {
		if err == sql.ErrNoRows {
			return nilSession, false, nil
		}

		return nilSession, false, err
	}

	return session, true, nil
}

func (r *DatabaseSessionsRepository) GetSessionWithUser(ctx context.Context, sessionId string) (models.Session, models.User, bool, error) {
	var user models.User
	var session models.Session
	row := r.db.GetConnection().QueryRowContext(ctx, "SELECT u.id, u.user_name, s.id, s.user_id, s.token_hash FROM sessions as s INNER JOIN users as u ON s.user_id = u.id WHERE s.id = $1", sessionId)
	if err := row.Scan(&user.Id, &user.UserName, &session.Id, &session.UserId, &session.TokenHash); err != nil {
		if err == sql.ErrNoRows {
			return nilSession, nilUser, false, nil
		}

		return nilSession, nilUser, false, err
	}

	return session, user, true, nil
}

func (r *DatabaseSessionsRepository) CreateUserSession(ctx context.Context, userId string) (models.Session, []byte, error) {
	session, token, err := models.NewUserSession(userId)
	if err != nil {
		return session, []byte{}, err
	}

	row := r.db.GetConnection().QueryRowContext(ctx, "INSERT INTO sessions (id, user_id, token_hash) VALUES (default, $1, $2) RETURNING id", session.UserId, session.TokenHash)
	if err := row.Scan(&session.Id); err != nil {
		return session, []byte{}, err
	}

	return session, token, nil
}

func (r *DatabaseSessionsRepository) DeleteSession(ctx context.Context, sessionId string) error {
	_, err := r.db.GetConnection().QueryContext(ctx, "DELETE FROM sessions WHERE id = $1", sessionId)
	if err != nil {
		return err
	}
	return nil
}

type NilSessionsRepository struct{}

func (r *NilSessionsRepository) GetSession(ctx context.Context, sessionId string) (models.Session, bool, error) {
	return nilSession, false, nil
}
func (r *NilSessionsRepository) GetSessionWithUser(ctx context.Context, sessionId string) (models.Session, models.User, bool, error) {
	return nilSession, nilUser, false, nil
}
func (r *NilSessionsRepository) CreateUserSession(ctx context.Context, userId string) (models.Session, []byte, error) {
	return nilSession, []byte{}, nil
}
func (r *NilSessionsRepository) DeleteSession(ctx context.Context, sessionId string) error {
	return nil
}
