package models

import (
	"time"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bookmark struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	ResourceID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	/* Relationships */
	User     User               `gorm:"foreignKey:UserID"`
	Resource WebCrawlerResource `gorm:"foreignKey:ResourceID"`
}

// BeforeCreate is a GORM hook that is triggered before a new Bookmark record is created in the database.
// It sets the CreatedAt field to the current time in the EAT (East Africa Time) timezone.
//
// Parameters:
//   - tx: The GORM database transaction object.
//
// Returns:
//   - err: An error if any issue occurs during the hook execution, otherwise nil.
func (b *Bookmark) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeUpdate is a GORM hook that is triggered before updating a Bookmark record.
// It updates the UpdatedAt field with the current time in the configured EAT timezone.
// This ensures that the record's last modification timestamp is accurate.
//
// Parameters:
//   - tx: The GORM database transaction.
//
// Returns:
//   - err: An error if any issues occur during the hook execution, otherwise nil.
func (b *Bookmark) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().In(config.EAT)
	return nil
}
