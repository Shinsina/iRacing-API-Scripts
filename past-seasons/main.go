package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type MinimalInitialResponse struct {
	Link string `json:"link"`
}

func main() {
	content, err := os.ReadFile("../cookie.txt")
	if err != nil {
		fmt.Println(1, err)
	}
	var cookies map[string]string
	err = json.Unmarshal(content, &cookies)
	if err != nil {
		fmt.Println(2, err)
	}
	raw_series_input, err := os.ReadFile("./distinct-series-ids-output.json")
	if err != nil {
		fmt.Println(3, err)
	}
	var unique_series_ids []int
	err = json.Unmarshal(raw_series_input, &unique_series_ids)
	if err != nil {
		fmt.Println(4, err)
	}
	channel := make(chan string, len(unique_series_ids))
	for i := 0; i < len(unique_series_ids); i++ {
		go func() {
			fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(unique_series_ids)) + " series")
			url := fmt.Sprintf("https://members-ng.iracing.com/data/series/past_seasons?series_id=%d", unique_series_ids[i])
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(5, err)
			}
			for key, value := range cookies {
				req.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client := &http.Client{}
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
			var initial_api_response MinimalInitialResponse
			err = json.Unmarshal(body, &initial_api_response)
			if err != nil {
				fmt.Println(8, err)
			}
			channel <- initial_api_response.Link
		}()
	}
	link_channel := make(chan string, len(unique_series_ids))
	for i := 0; i < len(unique_series_ids); i++ {
		fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(unique_series_ids)) + " links")
		url := <-channel
		go func() {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(9, err)
			}
			for key, value := range cookies {
				req.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client := &http.Client{}
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
	output_strings := make([]string, len(unique_series_ids))
	for i := 0; i < len(unique_series_ids); i++ {
		link_channel_response := <-link_channel
		output_strings[i] = link_channel_response
	}
	output_string := "["
	for i := 0; i < len(output_strings); i++ {
		output_string += (output_strings[i] + ",")
	}
	output_string += "]"
	file, err := os.Create("past-seasons-output.json")
	if err != nil {
		fmt.Println(12, err)
	}
	defer file.Close()
	_, err = file.Write([]byte(output_string))
	if err != nil {
		fmt.Println(13, err)
	}
}
