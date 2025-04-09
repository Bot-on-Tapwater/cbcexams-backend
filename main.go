/*
The main golang script/application
*/

package main

/*
Import config which contains database connection
Import database which runs database migrations
Import routes which contains the api endpoints routes
*/
import (
	"log"
	"os"

	"github.com/bot-on-tapwater/cbcexams-backend/config"
	"github.com/gin-gonic/gin"

	/*
	   Comment out "database" import after running first migration
	   There is an issue with the resources table
	   After importing data into the table trying
	   to run migrations will cause errors
	*/
	"github.com/bot-on-tapwater/cbcexams-backend/database"
	"github.com/bot-on-tapwater/cbcexams-backend/routes"
)

// main is the entry point of the application. It performs the following tasks:
//  1. Establishes a connection to the database using the ConnectDB function from the config package.
//  2. (Optional) Handles database migrations. Ensure that migrations are run only once or drop the
//     "resources" table before re-running migrations to avoid errors. Uncomment the
//     InitializeDatabase function call if migrations are required.
//  3. Sets up the HTTP router by passing the database connection to the SetupRouter function
//     from the routes package.
//  4. Starts the application on port 8080, making it ready to handle incoming HTTP requests.
func main() {
	/* Connect to the database */
	db := config.ConnectDB() // Ensure this function returns *gorm.DB

	/*
	   Only run migrations once
	   Or drop the "resources" table before running
	   migrations again to avoid errors
	*/
	/* Run database migrations */
	database.InitializeDatabase()

	/* Configure Gin */
	gin.SetMode(gin.ReleaseMode) // Switch to gin.DebugMode in development
	r := gin.Default()

	/* Register Routes */
	routes.AuthRoutes(r, db)
	routes.UsersRoutes(r, db)

	/* Specify the port to run your application and start server */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" /* Default port */
	}
	log.Printf("Server running on :%s", port)
	log.Fatal(r.Run(":" + port))
}
