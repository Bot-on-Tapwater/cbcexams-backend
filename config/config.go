package config

/*
gorm is our ORM
postgres allows us to connect to our postgres database
godotenv allows us to load environment variables from .env file
*/
import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
Global variable that represents a database connection
*/
var DB *gorm.DB

// ConnectDB establishes a connection to the PostgreSQL database using GORM.
// It loads environment variables from a .env file to construct the Data Source Name (DSN)
// and uses the DSN to open the database connection. If the connection is successful,
// it assigns the database instance to a global variable `DB` and returns the instance.
// If any error occurs during the process, the function logs the error and terminates the program.
//
// Returns:
//
//	*gorm.DB: A pointer to the GORM database instance.
//
// Environment Variables:
//
//	DB_HOST     - The hostname of the database server.
//	DB_USER     - The username for database authentication.
//	DB_PASSWORD - The password for database authentication.
//	DB_NAME     - The name of the database to connect to.
//	DB_PORT     - The port number on which the database server is running.
//
// Note:
//
//	Ensure that the .env file exists and contains the required environment variables
//	before calling this function.
func ConnectDB() *gorm.DB {
	/* Load environment variables */
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	/*
	   Construct Data Source Name (DSN) string
	   using environment variables
	*/
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	/*
	   Open the database connection using gorm.Open() method,
	   postgres.Open and the DSN
	*/
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	/*
	   Assign database instance to global variable DB
	*/
	DB = db
	fmt.Println("Database connected")
	return DB
}
