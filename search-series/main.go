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

type TokenReponse struct {
	Access_Token             string `json:"access_token"`
	Token_Type               string `json:"token_type"`
	Expires_In               int    `json:"expires_in"`
	Refresh_Token            string `json:"refresh_token"`
	Refresh_Token_Expires_In int    `json:"refresh_token_expires_in"`
	Scope                    string `json:"scope"`
}

type CustomerIDSeasonQuarterSeasonYearMapping struct {
	Customer_ID    int `json:"customer_id"`
	Season_Quarter int `json:"season_quarter"`
	Season_Year    int `json:"season_year"`
}

type iRacingAPIResponse struct {
	Data struct {
		Success    bool `json:"success"`
		Chunk_Info struct {
			Chunk_Size        int      `json:"chunk_size"`
			Num_Chunks        int      `json:"num_chunks"`
			Rows              int      `json:"rows"`
			Base_Download_Url string   `json:"base_download_url"`
			Chunk_File_Names  []string `json:"chunk_file_names"`
		} `json:"chunk_info"`
	} `json:"data"`
}

type InitialResponseGrouping struct {
	customer_id int
	outer_index int
	inner_index int
	response    iRacingAPIResponse
}

type ChunkResponseGrouping struct {
	customer_id int
	outer_index int
	response    string
}

func main() {

	content, err := os.ReadFile("../token.json")
	if err != nil {
		fmt.Println(err)
	}
	var token_response TokenReponse
	err = json.Unmarshal(content, &token_response)
	if err != nil {
		fmt.Println(err)
	}
	raw_array_of_arrays, err := os.ReadFile("./customer-id-season-quarter-season-year-mappings.json")
	if err != nil {
		fmt.Println(err)
	}
	var customer_id_season_quarter_season_year_array_of_arrays [][]CustomerIDSeasonQuarterSeasonYearMapping
	err = json.Unmarshal(raw_array_of_arrays, &customer_id_season_quarter_season_year_array_of_arrays)
	if err != nil {
		fmt.Println(err)
	}
	channel_count := 0
	for i := 0; i < len(customer_id_season_quarter_season_year_array_of_arrays); i++ {
		channel_count += len(customer_id_season_quarter_season_year_array_of_arrays[i])
	}
	channel := make(chan InitialResponseGrouping, channel_count)
	for i := 0; i < len(customer_id_season_quarter_season_year_array_of_arrays); i++ {
		for j := 0; j < len(customer_id_season_quarter_season_year_array_of_arrays[i]); j++ {
			mapping := customer_id_season_quarter_season_year_array_of_arrays[i][j]
			url := fmt.Sprintf("https://members-ng.iracing.com/data/results/search_series?season_quarter=%d&season_year=%d&cust_id=%d&official_only=true&event_types=5", mapping.Season_Quarter, mapping.Season_Year, mapping.Customer_ID)
			go func() {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					fmt.Println(err)
				}
				req.Header.Add("Authorization", token_response.Token_Type+" "+token_response.Access_Token)
				http_client := &http.Client{}
				// @todo Determine if this can be shortened even further
				sleep_time := 1000 * i
				time.Sleep(time.Duration(sleep_time) * time.Millisecond)
				resp, err := http_client.Do(req)
				if err != nil {
					fmt.Println(err)
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
				}
				var iracing_api_response iRacingAPIResponse
				err = json.Unmarshal(body, &iracing_api_response)
				if err != nil {
					fmt.Println(err, string(body))
				}
				var channel_output InitialResponseGrouping
				channel_output.customer_id = mapping.Customer_ID
				channel_output.outer_index = i
				channel_output.inner_index = j
				channel_output.response = iracing_api_response
				channel <- channel_output
			}()
		}
	}
	chunk_channel := make(chan ChunkResponseGrouping, channel_count)
	for i := 0; i < channel_count; i++ {
		channel_response := <-channel
		customer_id := channel_response.customer_id
		outer_index := channel_response.outer_index
		base_url := channel_response.response.Data.Chunk_Info.Base_Download_Url
		for i := 0; i < len(channel_response.response.Data.Chunk_Info.Chunk_File_Names); i++ {
			url := base_url + channel_response.response.Data.Chunk_Info.Chunk_File_Names[i]
			go func() {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					fmt.Println(err)
				}
				// Leaving this here in the event Signature and X-Amz-Algorithm are not query parameters at some point in time
				// req.Header.Add("Authorization", token_response.Token_Type+" "+token_response.Access_Token)
				http_client := &http.Client{}
				// @todo Determine if this can be shortened even further
				sleep_time := 1000 * i
				time.Sleep(time.Duration(sleep_time) * time.Millisecond)
				resp, err := http_client.Do(req)
				if err != nil {
					fmt.Println(err)
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
				}
				var chunk_channel_output ChunkResponseGrouping
				chunk_channel_output.outer_index = outer_index
				chunk_channel_output.customer_id = customer_id
				chunk_channel_output.response = string(body)
				chunk_channel <- chunk_channel_output
			}()
		}
	}
	output_strings := make([]string, len(customer_id_season_quarter_season_year_array_of_arrays))
	for i := 0; i < channel_count; i++ {
		chunk_channel_response := <-chunk_channel
		outer_index := chunk_channel_response.outer_index
		output_strings[outer_index] += (chunk_channel_response.response + ",")
	}
	for i := 0; i < len(customer_id_season_quarter_season_year_array_of_arrays); i++ {
		value := customer_id_season_quarter_season_year_array_of_arrays[i][0]
		customer_id := value.Customer_ID
		file_name := strconv.Itoa(customer_id) + "-search-series-output.json"
		file, err := os.Create(file_name)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		_, err = file.Write([]byte("[" + output_strings[i] + "]"))
		if err != nil {
			fmt.Println(err)
		}
	}
}
