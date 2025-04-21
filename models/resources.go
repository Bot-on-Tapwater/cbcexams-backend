package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// WebCrawlerResource represents a resource crawled from the web.
// It contains metadata and information about the resource, including its
// storage links, categorization, and extraction status.
//
// Fields:
// - ID: A unique identifier for the resource, generated as a UUID.
// - ParentURL: The URL of the parent resource from which this resource was crawled.
// - GoogleDriveDownloadLink: A unique Google Drive download link for the resource.
// - Name: The name of the resource.
// - RelativePath: The relative file path of the resource.
// - ParentDirectory: The directory containing the resource.
// - DjangoRelativePath: A unique relative path used in a Django application.
// - GoogleCloudStorageLink: A link to the resource stored in Google Cloud Storage.
// - CreatedAt: The timestamp when the resource was created.
// - Categories: A list of categories associated with the resource (requires github.com/lib/pq).
// - IsExtracted: A boolean indicating whether the resource's content has been extracted.
// - ExtractedContent: The extracted content of the resource, if available.
type WebCrawlerResource struct {
	ID                      uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ParentURL               string         `gorm:"type:text" json:"parent_url"`
	GoogleDriveDownloadLink string         `gorm:"type:text;uniqueIndex:web_crawler_resources_google_drive_download_link_key" json:"google_drive_download_link"`
	Name                    string         `gorm:"type:text" json:"name"`
	RelativePath            string         `gorm:"type:text" json:"relative_path"`
	ParentDirectory         string         `gorm:"type:text" json:"parent_directory"`
	DjangoRelativePath      string         `gorm:"type:text;uniqueIndex:web_crawler_resources_django_relative_path_key" json:"django_relative_path"`
	GoogleCloudStorageLink  string         `gorm:"type:text" json:"google_cloud_storage_link"`
	CreatedAt               time.Time      `gorm:"not null" json:"created_at"`
	Categories              pq.StringArray `gorm:"type:varchar(255)[]" json:"categories"` // Requires github.com/lib/pq
	IsExtracted             bool           `json:"is_extracted"`
	ExtractedContent        string         `gorm:"type:text" json:"extracted_content"`
}
