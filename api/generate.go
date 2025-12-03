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

const (
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	numberChars = "0123456789"
	symbolChars = "!@#$%^&*()_+~`|}{[]:;?><,./-="
)

// The ingredients: secure characters
func generateWithCharset(length int, charset string) (string, error) {
	if charset == "" {
		// Safety net: If ingredients are empty, use everything
		charset = upperChars + lowerChars + numberChars + symbolChars
	}

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	// 2. Parse Length
	lengthStr := r.URL.Query().Get("length")
	length := 12
	if lengthStr != "" {
		if l, err := strconv.Atoi(lengthStr); err == nil && l > 0 && l <= 100 {
			length = l
		}
	}

	// 3. Parse Filters (Dynamic Ingredients)
	// We build the 'validChars' string based on query params
	var validChars string

	if r.URL.Query().Get("upper") == "true" {
		validChars += upperChars
	}
	if r.URL.Query().Get("lower") == "true" {
		validChars += lowerChars
	}
	if r.URL.Query().Get("number") == "true" {
		validChars += numberChars
	}
	if r.URL.Query().Get("symbol") == "true" {
		validChars += symbolChars
	}

	// Note: 'Ambiguous' logic is usually handled on frontend for display,
	// but can be added here if needed. For now, we stick to basic filters.

	// 4. Generate
	password, err := generateWithCharset(length, validChars)
	if err != nil {
		http.Error(w, "Error generating password", http.StatusInternalServerError)
		return
	}

	// 5. Respond
	json.NewEncoder(w).Encode(PasswordResponse{
		Password: password,
		Length:   length,
	})
}
