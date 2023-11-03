package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jagottsicher/myGoWebApplication/internal/config"
	"github.com/jagottsicher/myGoWebApplication/internal/driver"
	"github.com/jagottsicher/myGoWebApplication/internal/handlers"
	"github.com/jagottsicher/myGoWebApplication/internal/helpers"
	"github.com/jagottsicher/myGoWebApplication/internal/models"
	"github.com/jagottsicher/myGoWebApplication/internal/render"
)

const portNumber = ":8080"
const versionNumber = "v1.0.176"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting E-Mail listener")
	listenForMail()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}

}

func run() (*driver.DB, error) {
	// Data to be available in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Bungalow{})
	gob.Register(models.BungalowRestriction{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags as arguments from the command line
	inProduction := flag.Bool("production", true, "Application is in production mode")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")
	version := flag.Bool("version", false, "Prints the version number")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Required flags missing - no user credentials for db access?")
		os.Exit(1)
	}

	if *version == true {
		fmt.Println(versionNumber)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// don't forget to change to true in Production!
	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connecting to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("No connection to database! Terminating ...")
	}
	log.Println("Successfully connected to database.")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)
	return db, nil
}
