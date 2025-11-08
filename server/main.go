package main

import (
	"log"
	"net/http"
)

func main() {
	config, err := MakeConfig()
	if err != nil {
		log.Fatalf("failed to init application: %s", err)
	}

	log.Println("starting http server")
	mux := MakeMux(config)
	// TODO: provide port as a config
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start http server: %s", err)
	}
}
