package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	a := App{}
	a.DB = NewPostgresConnection(os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	a.initializeRoutes()
	a.ensureTableExists()

	go func() {
		httpErr := http.ListenAndServe(":8080", a.Router)
		if httpErr != nil {
			log.Fatal(httpErr)
		}
		log.Println("server started")
	}()

	done := sigListener(func() {
		log.Println("closing db connection")
		a.DB.Close()
	})
	<-done
	log.Println("exiting server")
}

func sigListener(cleanup func()) chan bool {
	done := make(chan bool, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for _ = range c {
			cleanup()
			log.Println("exit signal received")
			done <- true
		}
	}()
	return done
}
