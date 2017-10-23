package main

import (
	"os"
	"log"
	"net/http"
	"github.com/ssharif6/Simplify/servers/handlers"
	"fmt"
)

func main() {
	portAddress := os.Getenv("ADDR")
	if len(portAddress) == 0 {
		// default port for HTTPS
		portAddress = ":443"
	}

	tlsKey := os.Getenv("TLSKEY")
	tlsCert := os.Getenv("TLSCERT")

	if len(tlsCert) == 0 || len(tlsKey) == 0 {
		log.Fatalf("TLS KEY OR TLS CERT NOT SET")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/simplify/text", handlers.TextHandler)
	mux.HandleFunc("/v1/simplify/image", handlers.ImageHandler)
	mux.HandleFunc("/v1/simplify/url", handlers.UrlHandler)

	fmt.Printf("Listening on port %d", portAddress)
	log.Fatal(http.ListenAndServeTLS(portAddress, tlsCert, tlsKey, mux))









	fmt.Println("HELLO DUBHACKS!")























	}