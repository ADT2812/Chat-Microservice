package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//////////////////////////////////////////////////
// 🔐 JWT CONFIG
//////////////////////////////////////////////////

var jwtKey = []byte("super-secret-key")

//////////////////////////////////////////////////
// 👤 USER STRUCT
//////////////////////////////////////////////////

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//////////////////////////////////////////////////
// 🗄 IN-MEMORY USER STORE
//////////////////////////////////////////////////

var users = make(map[string]string) // username -> hashed password

//////////////////////////////////////////////////
// ✅ ENABLE CORS FOR BROWSER
//////////////////////////////////////////////////

func enableCORS(w http.ResponseWriter) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

}

//////////////////////////////////////////////////
// 🔐 REGISTER API
//////////////////////////////////////////////////

func register(w http.ResponseWriter, r *http.Request) {

	enableCORS(w)

	// Handle Preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username & Password Required", http.StatusBadRequest)
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		http.Error(w, "Password Hash Error", http.StatusInternalServerError)
		return
	}

	users[user.Username] = string(hashedPassword)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User Registered Successfully",
	})
}

//////////////////////////////////////////////////
// 🔐 LOGIN API
//////////////////////////////////////////////////

func login(w http.ResponseWriter, r *http.Request) {

	enableCORS(w)

	// Handle Preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	storedPassword, exists := users[user.Username]
	if !exists {
		http.Error(w, "User Not Found", http.StatusUnauthorized)
		return
	}

	// Compare Password
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	// ✅ Create JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		},
	)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Token Generation Failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

//////////////////////////////////////////////////
// 🚀 MAIN
//////////////////////////////////////////////////

func main() {

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	log.Println("✅ Auth Service Running on :8001")

	log.Fatal(http.ListenAndServe(":8001", nil))
}
