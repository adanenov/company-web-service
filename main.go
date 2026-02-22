package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

// Краш-контрол для демонстрации restart
func crash(w http.ResponseWriter, r *http.Request) {
	log.Println("Container crash triggered")
	// Завершаем процесс полностью
	os.Exit(1)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	http.HandleFunc("/api", api)
	http.HandleFunc("/crash", crash)

	fmt.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
