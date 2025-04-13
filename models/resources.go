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
	ID                      uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ParentURL               string         `gorm:"type:text"`
	GoogleDriveDownloadLink string         `gorm:"type:text;uniqueIndex:web_crawler_resources_google_drive_download_link_key"`
	Name                    string         `gorm:"type:text"`
	RelativePath            string         `gorm:"type:text"`
	ParentDirectory         string         `gorm:"type:text"`
	DjangoRelativePath      string         `gorm:"type:text;uniqueIndex:web_crawler_resources_django_relative_path_key"`
	GoogleCloudStorageLink  string         `gorm:"type:text"`
	CreatedAt               time.Time      `gorm:"not null"`
	Categories              pq.StringArray `gorm:"type:varchar(255)[]"` // Requires github.com/lib/pq
	IsExtracted             bool
	ExtractedContent        string `gorm:"type:text"`
}
