package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"
)

type TokenReponse struct {
	Access_Token             string `json:"access_token"`
	Token_Type               string `json:"token_type"`
	Expires_In               int    `json:"expires_in"`
	Refresh_Token            string `json:"refresh_token"`
	Refresh_Token_Expires_In int    `json:"refresh_token_expires_in"`
	Scope                    string `json:"scope"`
}

type MinimalSubsession struct {
	Subsession_ID int `json:"subsession_id"`
}

type MinimalInitialResponse struct {
	Link string `json:"link"`
}

func main() {
	content, err := os.ReadFile("../token.json")
	if err != nil {
		fmt.Println(1, err)
	}
	var token_response TokenReponse
	err = json.Unmarshal(content, &token_response)
	if err != nil {
		fmt.Println(2, err)
	}
	raw_subsessions_input, err := os.ReadFile("./subsessions-input.json")
	if err != nil {
		fmt.Println(3, err)
	}
	var subsessions_minimal [][]MinimalSubsession
	err = json.Unmarshal(raw_subsessions_input, &subsessions_minimal)
	if err != nil {
		fmt.Println(4, err)
	}
	subsession_count := 0
	for i := 0; i < len(subsessions_minimal); i++ {
		subsession_count += len(subsessions_minimal[i])
	}
	var flattened_subsessions_ids []int
	for i := 0; i < len(subsessions_minimal); i++ {
		for j := 0; j < len(subsessions_minimal[i]); j++ {
			subsession_id := subsessions_minimal[i][j].Subsession_ID
			flattened_subsessions_ids = append(flattened_subsessions_ids, subsession_id)
		}
	}
	slices.Sort(flattened_subsessions_ids)
	// Highly suggested to process this in batches of 1000, if this value exceeds 1000 some requests may time out
	// To do this for the first 1000 do this
	// unique_subession_ids := slices.Compact(flattened_subsessions_ids)[:1000]
	// For the next subsequent 1000 do this
	// unique_subession_ids := slices.Compact(flattened_subsessions_ids)[1000:]
	// Adjust 1000 to 2000, 3000, etc. for each subsequent batch
	unique_subession_ids := slices.Compact(flattened_subsessions_ids)
	channel := make(chan string, len(unique_subession_ids))
	for i := 0; i < len(unique_subession_ids); i++ {
		go func() {
			fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(unique_subession_ids)) + " subsessions")
			url := fmt.Sprintf("https://members-ng.iracing.com/data/results/get?subsession_id=%d", unique_subession_ids[i])
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(5, err)
			}
			req.Header.Add("Authorization", token_response.Token_Type+" "+token_response.Access_Token)
			http_client := &http.Client{}
			// @todo Determine if this can be shortened even further
			sleep_time := 100 * i
			time.Sleep(time.Duration(sleep_time) * time.Millisecond)
			resp, err := http_client.Do(req)
			if err != nil {
				fmt.Println(6, err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(7, err)
			}
			// @todo Remove this once a stable sleep_time value has been determined
			// fmt.Println(resp)
			var initial_api_response MinimalInitialResponse
			err = json.Unmarshal(body, &initial_api_response)
			if err != nil {
				fmt.Println(8, err)
			}
			channel <- initial_api_response.Link
		}()
	}
	link_channel := make(chan string, len(unique_subession_ids))
	for i := 0; i < len(unique_subession_ids); i++ {
		fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(unique_subession_ids)) + " links")
		url := <-channel
		go func() {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(9, err)
			}
			// Leaving this here in the event Signature and X-Amz-Algorithm are not query parameters at some point in time
			// req.Header.Add("Authorization", token_response.Token_Type+" "+token_response.Access_Token)
			http_client := &http.Client{}
			// @todo Determine if this can be shortened even further
			sleep_time := 100 * i
			time.Sleep(time.Duration(sleep_time) * time.Millisecond)
			resp, err := http_client.Do(req)
			if err != nil {
				fmt.Println(10, err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(11, err)
			}
			link_channel <- string(body)
		}()
	}
	output_strings := make([]string, len(unique_subession_ids))
	for i := 0; i < len(unique_subession_ids); i++ {
		link_channel_response := <-link_channel
		output_strings[i] = link_channel_response
	}
	output_string := "["
	for i := 0; i < len(output_strings); i++ {
		output_string += (output_strings[i] + ",")
	}
	output_string += "]"
	file, err := os.Create("subsessions-output.json")
	if err != nil {
		fmt.Println(12, err)
	}
	defer file.Close()
	_, err = file.Write([]byte(output_string))
	if err != nil {
		fmt.Println(13, err)
	}
}
