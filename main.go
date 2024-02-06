package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
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
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
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

	middlewareList := []Middleware{RedisCache("")}

	dir := NewDavFileSystem(url, username, password)
	middleware := FileSystemMiddleware{
		FileSystem:  dir,
		Middlewares: middlewareList,
	}
	fMux := http.NewServeMux()
	fileServer := http.FileServer(middleware)
	fMux.Handle("/", fileServer)

	log.Printf("Listening...")
	log.Fatal(http.ListenAndServe(":8080", fMux))
}
