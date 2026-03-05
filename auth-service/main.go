package main

import (
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User Login Success")
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User Registered")
}

func main() {

	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)

	fmt.Println("Auth Service running on 8001")
	http.ListenAndServe(":8001", nil)

}
