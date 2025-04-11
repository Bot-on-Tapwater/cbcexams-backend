package models

import (
	"time"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

/* For students/parents seeking tutors */
type TutorRequest struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string         `gorm:"not null" json:"full_name"`
	Email          string         `gorm:"not null" json:"email"`
	Phone          string         `gorm:"not null" json:"phone_number"`
	Subjects       pq.StringArray `gorm:"type:text[];not null" json:"subjects"`
	EducationLevel pq.StringArray `gorm:"type:text[];not null" json:"education_level"`
	Location       string         `gorm:"not null" json:"location"`
	AvailableDays  pq.StringArray `gorm:"type:text[];not null" json:"available_days"`
	PreferredMode  string         `gorm:"not null" json:"preferred_mode"`
	AdditionalInfo string         `gorm:"type:text" json:"additional_info"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate is a GORM hook that is triggered before a new TutorRequest record
// is created in the database. It sets the CreatedAt field to the current time
// in the configured East Africa Time (EAT) timezone.
func (t *TutorRequest) BeforeCreate(tx *gorm.DB) (err error) {
	t.CreatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeUpdate is a GORM hook that is triggered before updating a TutorRequest record.
// It updates the UpdatedAt field with the current time in the configured EAT timezone.
// This ensures that the record's last modification timestamp is accurate.
func (t *TutorRequest) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now().In(config.EAT)
	return nil
}

/* For tutors offering services */
type TutorApplication struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string         `gorm:"not null" json:"full_name" form:"full_name"`
	Email          string         `gorm:"not null" json:"email" form:"email"`
	Phone          string         `gorm:"not null" json:"phone_number" form:"phone_number"`
	Subjects       pq.StringArray `gorm:"type:text[];not null" json:"subjects"`
	EducationLevel pq.StringArray `gorm:"type:text[];not null" json:"education_level"`
	Location       string         `gorm:"not null" json:"location" form:"location"`
	AvailableDays  pq.StringArray `gorm:"type:text[];not null" json:"available_days"`
	PreferredMode  string         `gorm:"not null" json:"preferred_mode" form:"preferred_mode"`
	ResumePath     string         `gorm:"size:255" json:"ResumePath"`
	AdditionalInfo string         `gorm:"type:text" json:"additional_info" form:"additional_info"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate is a GORM hook that is triggered before a new TutorApplication record
// is created in the database. It sets the CreatedAt field to the current time
// in the configured East Africa Time (EAT) timezone.
func (t *TutorApplication) BeforeCreate(tx *gorm.DB) (err error) {
	t.CreatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeUpdate is a GORM hook that is triggered before updating a TutorApplication record.
// It updates the UpdatedAt field with the current time in the configured EAT timezone.
// This ensures that the record's last modification timestamp is accurate.
func (t *TutorApplication) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now().In(config.EAT)
	return nil
}
