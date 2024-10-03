package main

import (
	"net/http"
	"os"

	middleware "davhttpd/middleware"

	log "github.com/sirupsen/logrus"
)

func init() {
	setupLogging()
}

func setupLogging() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "debug"
	}
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
}

func main() {
	url, ok := os.LookupEnv("URL")
	if !ok {
		panic("URL not set")
	}
	username, ok := os.LookupEnv("USER")
	if !ok {
		panic("USER not set")
	}
	password, ok := os.LookupEnv("PASSWORD")
	if !ok {
		panic("PASSWORD not set")

	}

	dir := NewDavFileSystem(url, username, password)

	middlewareList := []middleware.Middleware{
		&middleware.FileLogging{},
		&middleware.FileCache{},
	}

	mw := middleware.CreateChain(middlewareList, dir)

	fMux := http.NewServeMux()
	fileServer := http.FileServer(mw)
	fMux.Handle("/", middleware.HttpLogging(fileServer))

	log.Printf("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", fMux))
}
