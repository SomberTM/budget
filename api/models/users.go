package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id           string `json:"id"`
	UserName     string `json:"user_name"`
	PasswordHash []byte `json:"-"`
}

func NewUser() User {
	user := User{}
	return user
}

func (u *User) SetUserName(userName string) *User {
	u.UserName = userName
	return u
}

// Given a password as a raw string generates a bcrypt hash and
// updates the users hash given no errors occurred
func (u *User) SetPassword(password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return u, err
	}

	u.PasswordHash = hash
	return u, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	return err == nil
}
