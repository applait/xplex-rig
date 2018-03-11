package models

import (
	"time"

	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserAccount is used to store user information
type UserAccount struct {
	ID        uuid.UUID `sql:",pk,type:uuid"`
	Username  string    `sql:",unique,notnull"`
	Email     string    `sql:",unique,notnull"`
	Password  string    `json:"-"` // Ignore this field if this struct is marshalled to JSON
	IsActive  bool      `sql:",notnull,default:false"`
	CreatedAt time.Time `sql:",notnull"`
	UpdatedAt time.Time `sql:",notnull"`

	MultiStreams []*MultiStream // UserAccount hasMany MultiStream
}

// Find is a utility function to load user from DB based on values set in
// `User`. It returns `pg.ErrNoRows` if no user is found
func (u *UserAccount) Find(db *pg.DB) error {
	q := db.Model(u)
	if u.ID != uuid.Nil {
		q.Where("user_account.id = ?", u.ID)
	}
	if u.Username != "" {
		q.Where("user_account.username = ?", u.Username)
	}
	if u.Email != "" {
		q.Where("user_account.email = ?", u.Email)
	}
	return q.First()
}

// Insert current user in DB
func (u *UserAccount) Insert(db *pg.DB) (bool, error) {
	u.ID = uuid.NewV4()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	res, err := db.Model(u).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// Update current user in DB
func (u *UserAccount) Update(db *pg.DB) (bool, error) {
	u.UpdatedAt = time.Now()
	res, err := db.Model(u).Update()
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// UpdatePassword updates user's password
func (u *UserAccount) UpdatePassword(db *pg.DB, newPassword string) (bool, error) {
	err := u.SetPassword(newPassword)
	if err != nil {
		return false, err
	}
	return u.Update(db)
}

// UpdateEmail updates user's email
func (u *UserAccount) UpdateEmail(db *pg.DB, email string) (bool, error) {
	u.Email = email
	return u.Update(db)
}

// SetPassword hashes and stores user password
func (u *UserAccount) SetPassword(p string) error {
	password := []byte(p)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// MatchPassword matches plaintext password with stored hash password
func (u UserAccount) MatchPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return false
	}
	return true
}
