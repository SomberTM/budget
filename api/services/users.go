package services

import (
	"budget/api/models"
	"budget/api/repositories"
	"context"
	"errors"
)

type UsersService interface {
	Login(ctx context.Context, userName string, password string) (models.Session, []byte, error)
	Logout(ctx context.Context, sessionId string, sessionToken string) error
}

type SessionStrategyUsersService struct {
	usersRepository    repositories.UsersRepository
	sessionsRepository repositories.SessionsRepository
}

func NewSessionStrategyUsersService(usersRepository repositories.UsersRepository, sessionsRepository repositories.SessionsRepository) *SessionStrategyUsersService {
	return &SessionStrategyUsersService{usersRepository: usersRepository, sessionsRepository: sessionsRepository}
}
func (s *SessionStrategyUsersService) Login(ctx context.Context, userName string, password string) (models.Session, []byte, error) {
	user, exists, err := s.usersRepository.GetUserByUserName(ctx, userName)
	if err != nil {
		return models.Session{}, []byte{}, err
	}

	if !exists {
		user = models.NewUser()
		user.SetUserName(userName)
		user.SetPassword(password)

		user, err = s.usersRepository.CreateUser(ctx, user)
		if err != nil {
			return models.Session{}, []byte{}, err
		}
	}

	return s.sessionsRepository.CreateUserSession(ctx, user.Id)
}
func (s *SessionStrategyUsersService) Logout(ctx context.Context, sessionId string, sessionToken string) error {
	session, exists, err := s.sessionsRepository.GetSession(ctx, sessionId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("no session with id " + sessionId)
	}

	if !session.CheckToken([]byte(sessionToken)) {
		return errors.New("invalid session token")
	}

	err = s.sessionsRepository.DeleteSession(ctx, sessionId)
	if err != nil {
		return err
	}

	return nil
}

type NilUsersService struct{}

func (s *NilUsersService) Login(ctx context.Context, userName string, password string) (models.Session, []byte, error) {
	return models.Session{}, []byte{}, nil
}
func (s *NilUsersService) Logout(ctx context.Context, sessionId string, sessionToken string) error {
	return nil
}
