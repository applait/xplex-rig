package models

import (
	"time"

	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
)

// User is used to store user information
type User struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `sql:",unique,notnull"`
	Email     string `sql:",unique,notnull"`
	Password  string `sql:",notnull"`
	IsActive  bool   `sql:",notnull,default:false"`

	MultiStreams []*MultiStream // User hasMany MultiStream
}

// Find is a utility function to load user from DB based on values set in
// `User`. It returns `pg.ErrNoRows` if no user is found
func (u *User) Find(db *pg.DB) error {
	q := db.Model(u)
	if u.ID != 0 {
		q.Where("id = ?", u.ID)
	}
	if u.Username != "" {
		q.Where("username = ?", u.Username)
	}
	if u.Email != "" {
		q.Where("email = ?", u.Email)
	}
	return q.First()
}

// Insert current user in DB
func (u *User) Insert(db *pg.DB) (bool, error) {
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
func (u *User) Update(db *pg.DB) (bool, error) {
	u.UpdatedAt = time.Now()
	res, err := db.Model(&u).Update()
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

// UpdatePassword updates user's password
func (u *User) UpdatePassword(db *pg.DB, newPassword string) (bool, error) {
	err := u.SetPassword(newPassword)
	if err != nil {
		return false, err
	}
	return u.Update(db)
}

// UpdateEmail updates user's email
func (u *User) UpdateEmail(db *pg.DB, email string) (bool, error) {
	u.Email = email
	return u.Update(db)
}

// SetPassword hashes and stores user password
func (u *User) SetPassword(p string) error {
	password := []byte(p)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// MatchPassword matches plaintext password with stored hash password
func (u User) MatchPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return false
	}
	return true
}
