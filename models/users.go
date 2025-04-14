package models

import (
	"fmt"
	"time"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email                string     `gorm:"unique;not null" json:"email"`
	Password             string     `gorm:"not null" json:"password"` /* Hidden in JSON responses */
	FirstName            string     `gorm:"size:100" json:"first_name"`
	LastName             string     `gorm:"size:100" json:"last_name"`
	IsActive             bool       `gorm:"default:false" json:"is_active"`
	LastLogin            time.Time  `json:"last_login"`
	PasswordResetToken   string     `gorm:"size:255" json:"password_reset_token"`
	PasswordResetExpires time.Time  `json:"password_reset_expires"`
	Bookmarks            []Bookmark `gorm:"foreignKey:UserID"`
}

// BeforeLastLogin is a GORM hook that is triggered before updating the LastLogin field of a User.
// It sets the LastLogin field to the current time in the East Africa Time (EAT) timezone.
//
// Parameters:
//   - tx: The GORM database transaction object.
//
// Returns:
//   - err: An error if any issue occurs during the hook execution, otherwise nil.
func (u *User) BeforeLastLogin(tx gorm.DB) (err error) {
	u.LastLogin = time.Now().In(config.EAT)
	return nil
}

/* HashPassword encrypts the user's password before saving */
// HashPassword hashes the user's password using bcrypt with the default cost
// and updates the Password field with the hashed value. It returns an error
// if the hashing process fails.
func (u *User) HashPassword() error {
	fmt.Printf("Hashing the password: %s\n", u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

/* CheckPassword compares raw password with stored hash */
// CheckPassword compares the hashed password stored in the User struct
// with the provided plain-text password. It uses bcrypt to perform
// the comparison.
//
// Parameters:
//   - password: The plain-text password to be checked.
//
// Returns:
//   - An error if the passwords do not match or if there is an issue
//     during the comparison process. Returns nil if the passwords match.
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type UpdateProfileInput struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty" validate:"omitempty, email"`
}

func (u *User) UpdateProfile(input UpdateProfileInput) error {
	if input.FirstName != "" {
		u.FirstName = input.FirstName
	}

	if input.LastName != "" {
		u.LastName = input.LastName
	}

	if input.Email != "" && input.Email != u.Email {
		u.Email = input.Email
	}
	return nil
}
