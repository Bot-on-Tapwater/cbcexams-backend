package models

import (
	"time"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"not null" json:"full_Name"`
	Email string `gorm:"not null" json:"email"`
	Message string `gorm:"type:text;not null" json:"message"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate is a GORM hook that is triggered before a new Feedback record
// is created in the database. It sets the CreatedAt field to the current time
// in the East Africa Time (EAT) timezone.
func (f *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	f.CreatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeUpdate is a GORM hook that is triggered before updating a Feedback record.
// It updates the UpdatedAt field with the current time in the EAT (East Africa Time) timezone.
// 
// Parameters:
//   - tx: The GORM database transaction.
// 
// Returns:
//   - err: An error if any issues occur during the hook execution, otherwise nil.
func (f *Feedback) BeforeUpdate(tx *gorm.DB) (err error) {
	f.UpdatedAt = time.Now().In(config.EAT)
	return nil
}
