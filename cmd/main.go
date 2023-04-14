package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pipetail/sns_verification/internal/handlers"
	"github.com/pipetail/sns_verification/pkg/config"
)

func main() {
	cfg, err := config.ApplicationFromFlags()
	if err != nil {
		log.Fatalf("could not load the configuration: %s", err)
	}

	client := http.Client{}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.DefaultHandler(cfg, &client))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
