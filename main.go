package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
)

// 1. Define the "Menu" (What our JSON response looks like)
type PasswordResponse struct {
	Password string `json:"password"`
	Length   int    `json:"length"`
}

// The ingredients: secure characters
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"

// 2. The "Kitchen" (Logic to cook the password)
func generateSecurePassword(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		// Securely pick a random index
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

// 3. The "Waiter" (Handles the HTTP request)
func generateHandler(w http.ResponseWriter, r *http.Request) {
	// --- CORS HEADERS (Crucial!) ---
	// This tells the browser: "It's okay to accept requests from any website"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	
	// Parse the length from URL (e.g., ?length=20)
	lengthStr := r.URL.Query().Get("length")
	length := 12 // Default length
	
	if lengthStr != "" {
		if l, err := strconv.Atoi(lengthStr); err == nil && l > 0 && l <= 100 {
			length = l
		}
	}
	
	// Cook the password
	password, err := generateSecurePassword(length)
	if err != nil {
		http.Error(w, "Error generating password", http.StatusInternalServerError)
		return
	}
	
	// Serve the food (Send JSON)
	json.NewEncoder(w).Encode(PasswordResponse{
		Password: password,
		Length:   length,
	})
}

func main() {
	// Define the route
	http.HandleFunc("/generate", generateHandler)
	
	// Start the server
	fmt.Println("🚀 Backend is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
