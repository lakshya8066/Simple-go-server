package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func main() {

	PORT := 1337

	// Define the function to handle incoming requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			// Read the request body
			jsonstr, err := ioutil.ReadAll(r.Body)

			var data map[string]interface{}
			err = json.Unmarshal([]byte(jsonstr), &data)
			if err != nil {
				panic(err)
			}

			// Process the request body
			// fmt.Printf("Processing POST request with body:\n%s\n", body)

			// Process the POST request and generate the JWT token
			email := data["email"].(string)
			fmt.Printf("This is the email: %s\n", email)
			claims := &Claims{
				Email: email,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte("mysecretkey"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Set the JWT token in the response header
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

			//Sending a POST request to localhost:10998
			// Create a new POST request with the string as the body
			req, err := http.NewRequest("POST", "http://localhost:10998/token", bytes.NewBuffer([]byte(tokenString)))
			if err != nil {
				panic(err)
			}
			// Set the content type header to text/plain
			req.Header.Set("Content-Type", "text/plain")

			// Send the request and get the response
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			// Print the response status code and body
			fmt.Printf("Response Status Code: %d\n", resp.StatusCode)

			// Write the response
			// Process the POST request
			fmt.Fprintf(w, "Processing POST request hahaha\n")

		} else {
			fmt.Fprintf(w, "Invalid request method\n")
		}
	})

	// Start the server and listen on port 8080
	fmt.Printf("Starting server on port %d...\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
