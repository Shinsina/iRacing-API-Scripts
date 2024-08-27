package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

/**
* Needs to be executed via auth.sh
* Additionally IRACING_USERNAME and IRACING_PASSWORD must be set in a .env file at root
**/

func main() {
	initial_hash := sha256.Sum256([]byte(os.Getenv("IRACING_PASSWORD") + os.Getenv("IRACING_USERNAME")))
	hash_in_base_64 := base64.StdEncoding.EncodeToString(initial_hash[:])
	values := map[string]string{"email": os.Getenv("IRACING_USERNAME"), "password": hash_in_base_64}
	json_data, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err)
	}
	auth_response, err := http.Post("https://members-ng.iracing.com/auth", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
	}
	auth_json := make(map[string]string)
	cookies := auth_response.Cookies()
	for i := 0; i < len(cookies); i++ {
		_, exists := auth_json[cookies[i].Name]
		if !exists {
			new_value := cookies[i].Value
			auth_json[cookies[i].Name] = new_value
		}
	}
	cookie_as_bytes, err := json.Marshal(auth_json)
	if err != nil {
		fmt.Println(err)
	}
	file, err := os.Create("cookie.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.Write(cookie_as_bytes)
	if err != nil {
		fmt.Println(err)
	}
}
