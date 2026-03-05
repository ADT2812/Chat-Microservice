package main

import (
	"fmt"
	"net/http"
)

func notify(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Notification sent!")
}

func main() {

	http.HandleFunc("/notify", notify)

	fmt.Println("Notification Service running on 8004")

	http.ListenAndServe(":8004", nil)
}
