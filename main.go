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
	"github.com/gin-contrib/secure"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
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

	/* Initialize the EAT timezone */
	config.InitTimezone()

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

	/* Configure rate limiting */
	/* TODO: Reset after testing */
	rate, _ := limiter.NewRateFromFormatted("3000000-M") /* 300 requests per minute */
	store := memory.NewStore()
	middleware := ginlimiter.NewMiddleware(limiter.New(store, rate))

	/* Apply rate limiting middleware */
	r.Use(middleware)

	/* CORS configuration */
	r.Use(secure.New(secure.Config{
		/* Specifies the list of allowed hostnames. Requests with a Host header not in this list will be rejected. */
		// AllowedHosts: []string{
		// 	"vercel.app",
		// 	"*vercel.app",
		// 	"cbcexams.com",
		// 	"*cbcexams.com",
		// 	"localhost",
		// 	"localhost:8080",
		// 	"127.0.0.1",
		// 	"127.0.0.1:8080",
		// 	"",
		// },

		/* Redirects all HTTP requests to HTTPS if set to true. Disabled here for local development. */
		SSLRedirect: false,

		/* Specifies the duration (in seconds) for which the browser should remember that the site must only be accessed using HTTPS. */
		/* 31536000 seconds = 1 year. */
		STSSeconds: 31536000,

		/* If true, applies the Strict-Transport-Security (HSTS) policy to all subdomains as well. */
		STSIncludeSubdomains: true,

		/* Prevents the site from being displayed in an iframe to protect against clickjacking attacks. */
		FrameDeny: true,

		/* Prevents the browser from trying to guess the MIME type of a file, reducing the risk of MIME-based attacks. */
		ContentTypeNosniff: true,

		/* Enables the X-XSS-Protection header to prevent some types of cross-site scripting (XSS) attacks. */
		BrowserXssFilter: true,

		/* Prevents Internet Explorer from executing downloads in the site's context, reducing the risk of drive-by downloads. */
		IENoOpen: true,

		/* Specifies the Referrer-Policy header, controlling how much referrer information is sent with requests. */
		/* "strict-origin-when-cross-origin" sends the full URL for same-origin requests but only the origin for cross-origin requests. */
		ReferrerPolicy: "strict-origin-when-cross-origin",

		/* Specifies the Content-Security-Policy (CSP) header, which restricts the sources from which content can be loaded. */
		/* "default-src 'self'" allows content to be loaded only from the same origin. */
		ContentSecurityPolicy: "default-src 'self'",
	}))

	/* Register Routes */
	routes.AuthRoutes(r, db)
	routes.UsersRoutes(r, db)
	routes.CategoriesRoutes(r)
	routes.TutoringRoutes(r, db)
	routes.JobRoutes(r, db)
	routes.WebDevRoutes(r, db)
	routes.FeedbackRoutes(r, db)
	routes.BookmarkRoutes(r, db)
	routes.PaymentRoutes(r)
	routes.ResourceRoutes(r, db)

	/* Print all registered routes */
	for _, route := range r.Routes() {
		log.Printf("Method: %s | Path: %s", route.Method, route.Path)
	}

	/* Specify the port to run your application and start server */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" /* Default port */
	}
	log.Printf("Server running on :%s", port)
	log.Fatal(r.Run(":" + port))
}
