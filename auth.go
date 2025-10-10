package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

/**
* Needs to be executed via auth.sh
* Additionally IRACING_USERNAME, IRACING_PASSWORD, IRACING_CLIENT_ID and IRACING_CLIENT_SECRET must be set in a .env file at root
**/

func main() {
	iracing_username := os.Getenv("IRACING_USERNAME")
	iracing_password := os.Getenv("IRACING_PASSWORD")
	initial_username_password_hash := sha256.Sum256([]byte(iracing_password + iracing_username))
	username_password_hash_in_base_64 := base64.StdEncoding.EncodeToString(initial_username_password_hash[:])
	iracing_client_id := os.Getenv("IRACING_CLIENT_ID")
	iracing_client_secret := os.Getenv("IRACING_CLIENT_SECRET")
	initial_client_hash := sha256.Sum256([]byte(iracing_client_secret + iracing_client_id))
	client_hash_in_base_64 := base64.StdEncoding.EncodeToString(initial_client_hash[:])
	request_body := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s&scope=%s", url.QueryEscape("password_limited"), url.QueryEscape(iracing_client_id), url.QueryEscape(client_hash_in_base_64), url.QueryEscape(iracing_username), url.QueryEscape(username_password_hash_in_base_64), url.QueryEscape("iracing.auth"))
	auth_response, err := http.Post("https://oauth.iracing.com/oauth2/token", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(request_body)))
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(auth_response.Body)
	file, err := os.Create("token.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		fmt.Println(err)
	}
}
