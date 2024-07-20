package models

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	TokenHash []byte `json:"token_hash"`
}

func generateSessionToken(length int) ([]byte, error) {
	token := make([]byte, length)

	_, err := rand.Read(token)
	if err != nil {
		return []byte{}, err
	}

	base64Token := base64.URLEncoding.EncodeToString(token)
	return []byte(base64Token), nil
}

func HashSessionToken(token []byte) ([]byte, error) {
	// Figure out how to properly hash the session token
	// return token, nil

	hash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hash, nil
}

func NewSession() Session {
	session := Session{}
	return session
}

func NewUserSession(userId string) (Session, []byte, error) {
	session := NewSession()

	token, err := generateSessionToken(32)
	if err != nil {
		return session, []byte{}, err
	}

	hash, err := HashSessionToken(token)
	if err != nil {
		return session, []byte{}, err
	}

	session.TokenHash = hash
	session.UserId = userId
	return session, token, nil
}

func (s *Session) SetSessionCookie(c *gin.Context, token []byte) {
	c.SetCookie("session_id", s.Id, 3600, "/", "", true, true)
	c.SetCookie("session_token", string(token), 3600, "/", "", true, true)
}

func (s *Session) CheckToken(token []byte) bool {
	err := bcrypt.CompareHashAndPassword(s.TokenHash, token)
	return err == nil
}

func ClearSessionCookie(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", true, true)
	c.SetCookie("session_token", "", -1, "/", "", true, true)
}
