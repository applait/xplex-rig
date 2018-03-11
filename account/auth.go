package account

import (
	"errors"
)

// AuthLocal authenticates using username and password and returns a JWT for user if successful
func AuthLocal(username string, password string) (string, error) {
	var t string
	u, err := GetUserByUsername(username)
	if err != nil {
		return t, err
	}
	err = ValidatePassword(u, password)
	if err != nil {
		return t, errors.New("Invalid password")
	}
	t, err = newUserToken(u.ID, u.Username)
	if err != nil {
		return t, errors.New("Unable to authenticate user")
	}
	return t, nil
}
