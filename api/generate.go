package handler

import (
	"crypto/rand"
	"encoding/json"
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
func Handler(w http.ResponseWriter, r *http.Request) {
	// 1. CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	// 2. Parse Query
	lengthStr := r.URL.Query().Get("length")
	length := 12
	if lengthStr != "" {
		if l, err := strconv.Atoi(lengthStr); err == nil && l > 0 && l <= 100 {
			length = l
		}
	}

	// 3. Generate
	password, err := generateSecurePassword(length)
	if err != nil {
		http.Error(w, "Error generating password", http.StatusInternalServerError)
		return
	}

	// 4. Respond
	json.NewEncoder(w).Encode(PasswordResponse{
		Password: password,
		Length:   length,
	})
}
