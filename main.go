package main

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"agones.dev/agones/pkg/util/runtime"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"time"
)

var agonesClient versioned.Interface
var db *gorm.DB

func main() {
	config, err := rest.InClusterConfig()
	logger := runtime.NewLoggerWithSource("main")

	if err != nil {
		logger.WithError(err).Fatal("Could not create in-cluster config")
	}

	agonesClient, err = versioned.NewForConfig(config)
	if err != nil {
		logger.WithError(err).Fatal("Could not create agones api clientset")
	}

	db, err = gorm.Open(
		postgres.Open(fmt.Sprintf(
			"host=postgres user=%s password=%s database=%s port=5432 sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"))),
		&gorm.Config{})
	if err != nil {
		logger.WithError(err).Fatal("Could not connect to database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Score{})
	if err != nil {
		logger.WithError(err).Fatal("Could not migrate database")
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello world"))

		if err != nil {
			logger.WithError(err).Warn(err.Error())
			return
		}
	})

	r.Get("/api/rooms", createRoom)
	r.Get("/api/rooms/{room}", getRoom)

	r.Post("/api/scores", addScore)
	r.Get("/api/scores", getLeaderboard)
	r.Delete("/api/scores/{id}", deleteScore)

	r.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/leaderboard.html")
	})

	fs := http.FileServer(http.Dir("./client/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("404: Even the pigeons couldn't find this one"))

		if err != nil {
			logger.WithError(err).Warn(err.Error())
			return
		}
	})

	// TODO: Add TLS when it is actually hosted on AWS

	srv := &http.Server{
		Addr:           ":3000",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Fatal(srv.ListenAndServe())
}
