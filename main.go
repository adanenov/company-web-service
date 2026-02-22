package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "IT Company Website")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "About our company")
}

func api(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"status": "ok"}`)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	http.HandleFunc("/api", api)

	fmt.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
