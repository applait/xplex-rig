package account

import (
	"github.com/applait/xplex-rig/common"
	"github.com/satori/go.uuid"
)

// GetUserByID fetches a user account by the account ID
func GetUserByID(userID uuid.UUID) (common.UserAccount, error) {
	var u common.UserAccount
	query := `
    select
      id, username, email, password, is_active, created_at, updated_at
    from user_accounts
      where id = $1;
  `
	err := common.DB.QueryRow(query, userID).Scan(
		&u.ID, &u.Username, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetUserByUsername fetches a user account by the account ID
func GetUserByUsername(username string) (common.UserAccount, error) {
	var u common.UserAccount
	query := `
    select
      id, username, email, password, is_active, created_at, updated_at
    from user_accounts
      where username = $1;
  `
	err := common.DB.QueryRow(query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetUserByEmail fetches a user account by their email
func GetUserByEmail(email string) (common.UserAccount, error) {
	var u common.UserAccount
	query := `
    select
      id, username, email, password, is_active, created_at, updated_at
    from user_accounts
      where username = $1;
  `
	err := common.DB.QueryRow(query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// CreateUser creates a new user account
func CreateUser(u common.UserAccount) error {
	createQuery := `
    insert into user_accounts
      (id, username, email, password, is_active, created_at, updated_at)
    values
      ($1, $2, $3, $4, $5, now(), now())
    returning created_at, updated_at;
  `
	u.ID = uuid.NewV4()
	passwd, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	_, err = common.DB.Exec(createQuery,
		&u.ID,
		&u.Username,
		&u.Email,
		&passwd,
		false,
	)
	if err != nil {
		return err
	}
	// TODO: Send email notification
	return nil
}
