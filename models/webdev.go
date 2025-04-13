package models

import (
	"time"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebDevRequest struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ClientType       string    `gorm:"not-null" json:"client_type"`
	OrganizationName string    `json:"organization_name"`
	ContactName      string    `gorm:"not null" json:"full_name"`
	ContactEmail     string    `gorm:"not null" json:"email"`
	ContactPhone     string    `gorm:"not null" json:"phone_number"`
	ProjectType      string    `gorm:"not null" json:"project_type"`
	BudgetRange      string    `gorm:"not null" json:"budget_range"`
	ProjectDetails   string    `gorm:"type:text;not null" json:"project_details"`
	Status           string    `gorm:"default:'pending'" json:"status"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate is a GORM hook method that is triggered before a new WebDevRequest
// record is created in the database. It sets the CreatedAt field to the current
// time in the EAT (East Africa Time) timezone.
//
// Parameters:
// - tx: The GORM database transaction object.
//
// Returns:
// - err: An error if any issues occur during the hook execution, otherwise nil.
func (wdr *WebDevRequest) BeforeCreate(tx *gorm.DB) (err error) {
	wdr.CreatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeUpdate is a GORM hook that is triggered before updating a WebDevRequest record.
// It sets the UpdatedAt field to the current time in the EAT (East Africa Time) timezone.
// This ensures that the record's last updated timestamp is accurate.
//
// Parameters:
//   - tx: The GORM database transaction object.
//
// Returns:
//   - err: An error if any issue occurs during the hook execution, otherwise nil.
func (wdr *WebDevRequest) BeforeUpdate(tx *gorm.DB) (err error) {
	wdr.UpdatedAt = time.Now().In(config.EAT)
	return nil
}
