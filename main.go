/*
Demo logs server

Receive and view JSON data on your very own log panel.

Created to demonstrate the risk of infostealers
*/
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed static/*
var content embed.FS

type Rdata map[string]any

type Payload struct {
	Timestamp int64  `json:"timestamp"`
	UserAgent string `json:"useragent"`
	Data      Rdata  `json:"data"`
}

func newPayload(data map[string]any, useragent string) *Payload {
	return &Payload{
		Timestamp: time.Now().UnixMilli(),
		UserAgent: useragent,
		Data:      data,
	}
}

func main() {
	authHeader := os.Getenv("AUTH_HEADER")
	if authHeader == "" {
		log.Fatal("You must be set the AUTH_HEADER environment variable")
	}
	recieved := make([]*Payload, 0)
	// Send data (JSON) to the server
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("%s: %s: rejected method %s", r.URL, r.RemoteAddr, r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("NOT allowed\n"))

			return
		}
		auth := r.Header.Get("authorization")
		if auth == "" || auth != authHeader {
			log.Printf("%s: %s: no authorisation", r.URL, r.RemoteAddr)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden\n"))
			return
		}
		var j Rdata
		err := json.NewDecoder(r.Body).Decode(&j)
		if err != nil {
			log.Print(fmt.Errorf("reading body: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error\n"))
			return
		}
		ua := r.Header.Get("user-agent")
		log.Printf("%s: %s: accepted payload", r.URL, r.RemoteAddr)

		recieved = append(recieved, newPayload(j, ua))
		w.Write([]byte("received\n"))
	})
	// View collected data
	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("%s: %s: rejected method %s", r.URL, r.RemoteAddr, r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		auth := r.Header.Get("authorization")
		if auth == "" || auth != authHeader {
			log.Printf("%s: %s: no authorisation", r.URL, r.RemoteAddr)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(recieved)
	})
	// Clear collected data
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("%s: %s: rejected method %s", r.URL, r.RemoteAddr, r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		auth := r.Header.Get("authorization")
		if auth == "" || auth != authHeader {
			log.Printf("%s: %s: no authorisation", r.URL, r.RemoteAddr)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		recieved = make([]*Payload, 0)
		w.Write([]byte("cleared\n"))

	})
	contentStatic, err := fs.Sub(content, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(contentStatic)))
	log.Println("listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
