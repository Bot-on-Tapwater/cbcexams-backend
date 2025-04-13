package models

import (
	"time"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

/* School job listings */
type SchoolJobListing struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SchoolName          string         `gorm:"not null" json:"school_name"`
	ContactEmail        string         `gorm:"not null" json:"email"`
	ContactPhone        string         `gorm:"not null" json:"phone_number"`
	Position            string         `gorm:"not null" json:"position"`
	Subjects            pq.StringArray `gorm:"type:text[];not null" json:"subjects"`
	EducationLevel      pq.StringArray `gorm:"type:text[];not null" json:"education_level"`
	Location            string         `gorm:"not null" json:"location"`
	EmploymentType      string         `gorm:"not null" json:"employment_type"`
	ApplicationDeadline time.Time      `gorm:"not null" json:"application_deadline"`
	Description         string         `gorm:"type:text;not null" json:"description"`
	Requirements        string         `gorm:"type:text;not null" json:"requirements"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate is a GORM hook that is triggered before a new SchoolJobListing
// record is created in the database. It sets the CreatedAt field to the current
// time in the East Africa Time (EAT) timezone.
func (sjl *SchoolJobListing) BeforeCreate(tx *gorm.DB) (err error) {
	sjl.CreatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeUpdate is a GORM hook that is triggered before updating a SchoolJobListing record.
// It updates the UpdatedAt field with the current time in the configured EAT timezone.
//
// Parameters:
//   - tx: The GORM database transaction object.
//
// Returns:
//   - err: An error if any issues occur during the hook execution, otherwise nil.
func (sjl *SchoolJobListing) BeforeUpdate(tx *gorm.DB) (err error) {
	sjl.UpdatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeApplicationDeadline is a GORM hook method that is triggered before
// saving or updating a SchoolJobListing record. It sets the ApplicationDeadline
// field to the current time in the East Africa Time (EAT) timezone.
//
// Parameters:
// - tx: The GORM database transaction object.
//
// Returns:
// - err: An error if any issues occur during the execution of the hook, otherwise nil.
func (sjl *SchoolJobListing) BeforeApplicationDeadline(tx *gorm.DB) (err error) {
	sjl.ApplicationDeadline = time.Now().In(config.EAT)
	return nil
}

/* TeacherJobProfile */
type TeacherJobProfile struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FullName          string         `gorm:"not null" json:"full_name" form:"full_name"`
	Email             string         `gorm:"not null" json:"email" form:"email"`
	Phone             string         `gorm:"not null" json:"phone_number" form:"phone_number"`
	Subjects          pq.StringArray `gorm:"type:text[];not null" json:"subjects" form:"subjects"`
	EducationLevel    pq.StringArray `gorm:"type:text[];not null" json:"education_level" form:"education_level"`
	ExperienceYears   int            `gorm:"not null" json:"experience_years" form:"experience_years"`
	Location          string         `gorm:"not null" json:"location" form:"location"`
	WillingToRelocate bool           `gorm:"not null" json:"willing_to_relocate" form:"willing_to_relocate"`
	ResumePath        string         `gorm:"not null;size:255" json:"ResumePath"`
	IsActive          bool           `gorm:"default:true" json:"is_active" form:"is_active"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate is a GORM hook that is triggered before a new TeacherJobProfile
// record is created in the database. It sets the CreatedAt field to the current
// time in the East Africa Time (EAT) timezone.
func (tjp *TeacherJobProfile) BeforeCreate(tx *gorm.DB) (err error) {
	tjp.CreatedAt = time.Now().In(config.EAT)
	return nil
}

// BeforeUpdate is a GORM hook that is triggered before updating a TeacherJobProfile record.
// It updates the UpdatedAt field with the current time in the configured EAT timezone.
// This ensures that the record's last modification timestamp is accurate.
func (tjp *TeacherJobProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	tjp.UpdatedAt = time.Now().In(config.EAT)
	return nil
}
