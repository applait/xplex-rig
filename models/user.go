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

	RTMPStreams []*RTMPStream // User hasMany RTMPStream
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
func (u *User) Insert(db *pg.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return db.Insert(u)
}

// Update current user in DB
func (u *User) Update(db *pg.DB) error {
	u.UpdatedAt = time.Now()
	return db.Update(u)
}

// UpdatePassword updates user's password
func (u *User) UpdatePassword(db *pg.DB, newPassword string) error {
	err := u.SetPassword(newPassword)
	if err != nil {
		return err
	}
	return u.Update(db)
}

// UpdateEmail updates user's email
func (u *User) UpdateEmail(db *pg.DB, email string) error {
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
