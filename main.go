package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"agones.dev/agones/pkg/client/clientset/versioned"
	//"agones.dev/agones/pkg/util/runtime"
	//"k8s.io/client-go/rest"
)

var agonesClient versioned.Interface

func main() {
	//config, err := rest.InClusterConfig()
	//logger := runtime.NewLoggerWithSource("main")
	//
	//if err != nil {
	//	logger.WithError(err).Fatal("Could not create in-cluster config")
	//}
	//
	//agonesClient, err = versioned.NewForConfig(config)
	//if err != nil {
	//	logger.WithError(err).Fatal("Could not create agones api clientset")
	//}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Get("/rooms/", createRoom)
	r.Get("/rooms/{room}/", getRoom)

	err := http.ListenAndServe(":3000", r)

	if err != nil {
		print(err.Error())
		return
	}
}
