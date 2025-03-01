package main

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"agones.dev/agones/pkg/util/runtime"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"k8s.io/client-go/rest"
	"net/http"
	"time"
)

var agonesClient versioned.Interface

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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello world"))

		if err != nil {
			print(err.Error())
			return
		}
	})

	r.Get("/rooms/", createRoom)
	r.Get("/rooms/{room}/", getRoom)

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
