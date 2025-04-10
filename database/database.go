package database

/*
config allows us to connect to the database
models contains our database tables to be migrated
*/
import (
	"fmt"
	"log"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/bot-on-tapwater/cbcexams-backend/models"
)

// InitializeDatabase connects to the database and runs migrations.
// InitializeDatabase sets up the database connection and performs
// the necessary migrations for the application. It ensures that
// the database schema is up-to-date by applying migrations for
// the specified models. If the migration process fails, the
// function logs a fatal error and terminates the application.
func InitializeDatabase() {
	/* Initialise database connection */
	db := config.DB

	/* Run migrations */
	fmt.Println("Running database migrations...")
	err := db.AutoMigrate(&models.User{}) // Add more models here
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migrated successfully!")
}
