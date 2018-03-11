package account

import (
	"errors"
	"log"

	"github.com/applait/xplex-rig/common"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash generates a hash for a given password string
func generatePasswordHash(passwd string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return string(hashByte), err
	}

	return string(hashByte), nil
}

// VerifyPasswordHash verifies a hash against its password
func verifyPasswordHash(passwd string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	if err != nil {
		return err
	}

	return nil
}

// ValidatePassword verifies a user's password against a given password
func ValidatePassword(u common.UserAccount, password string) error {
	return verifyPasswordHash(password, u.Password)
}

// ChangePassword changes account password to new password assuming old password matches
func ChangePassword(userID uuid.UUID, oldPasswd string, newPasswd string) error {
	u, err := GetUserByID(userID)
	if err != nil {
		return err
	}
	err = verifyPasswordHash(oldPasswd, u.Password)
	if err != nil {
		return err
	}
	hashed, err := generatePasswordHash(newPasswd)
	if err != nil {
		return err
	}
	query := `
    update user_accounts
      set password = $1, updated_at = now()
    where id = $2;
  `
	res, err := common.DB.Exec(query, hashed, userID)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil || affected != 1 {
		return errors.New("Unable to update account with new password")
	}
	return nil
}
