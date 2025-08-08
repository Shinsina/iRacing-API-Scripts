package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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

	content, err := os.ReadFile("../cookie.txt")
	if err != nil {
		fmt.Println(err)
	}
	var cookies map[string]string
	err = json.Unmarshal(content, &cookies)
	if err != nil {
		fmt.Println(err)
	}
	raw_array_of_arrays, err := os.ReadFile("./customer-id-season-quarter-season-year-mappings.json")
	if err != nil {
		fmt.Println(err)
	}
	var customer_id_season_quarter_season_year_array []string
	err = json.Unmarshal(raw_array_of_arrays, &customer_id_season_quarter_season_year_array)
	if err != nil {
		fmt.Println(err)
	}
	channel_count := len(customer_id_season_quarter_season_year_array)
	channel := make(chan InitialResponseGrouping, channel_count)
	for i := 0; i < len(customer_id_season_quarter_season_year_array); i++ {
		mapping := customer_id_season_quarter_season_year_array[i]
		split_mapping := strings.Split(mapping, "_")
		url := fmt.Sprintf("https://members-ng.iracing.com/data/results/search_series?season_quarter=%s&season_year=%s&cust_id=%s&official_only=true&event_types=5", split_mapping[0], split_mapping[1], split_mapping[2])
		go func() {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}
			for key, value := range cookies {
				req.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client := &http.Client{}
			// @todo Determine if this can be shortened even further
			sleep_time := 1000 * i
			time.Sleep(time.Duration(sleep_time) * time.Millisecond)
			fmt.Println(fmt.Sprintf("Initial request for season quarter %s, season year %s, customer ID %s (%d of %d)", split_mapping[0], split_mapping[1], split_mapping[2], i+1, len(customer_id_season_quarter_season_year_array)))
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
			channel_output.customer_id, _ = strconv.Atoi(split_mapping[2])
			channel_output.outer_index = i
			channel_output.response = iracing_api_response
			channel <- channel_output
		}()
	}
	chunk_channel := make(chan ChunkResponseGrouping, channel_count)
	for i := 0; i < channel_count; i++ {
		fmt.Println(fmt.Sprintf("Consuming channel %d of %d", i+1, channel_count))
		channel_response := <-channel
		customer_id := channel_response.customer_id
		outer_index := channel_response.outer_index
		base_url := channel_response.response.Data.Chunk_Info.Base_Download_Url
		for i := 0; i < len(channel_response.response.Data.Chunk_Info.Chunk_File_Names); i++ {
			fmt.Println(fmt.Sprintf("Processing chunks from channel %d of %d", i+1, len(channel_response.response.Data.Chunk_Info.Chunk_File_Names)))
			url := base_url + channel_response.response.Data.Chunk_Info.Chunk_File_Names[i]
			go func() {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					fmt.Println(err)
				}
				for key, value := range cookies {
					req.AddCookie(&http.Cookie{Name: key, Value: value})
				}
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
	output_strings := make([]string, len(customer_id_season_quarter_season_year_array))
	for i := 0; i < channel_count; i++ {
		fmt.Println(fmt.Sprintf("Handling processed chunks from channel %d of %d", i+1, channel_count))
		chunk_channel_response := <-chunk_channel
		outer_index := chunk_channel_response.outer_index
		output_strings[outer_index] += (chunk_channel_response.response + ",")
	}
	cust_id_to_subsession_map := make(map[string][]string)
	for i := 0; i < len(customer_id_season_quarter_season_year_array); i++ {
		fmt.Println(fmt.Sprintf("Generating output for %d of %d channels", i+1, len(customer_id_season_quarter_season_year_array)))
		value := customer_id_season_quarter_season_year_array[i]
		split_value := strings.Split(value, "_")
		customer_id := split_value[2]
		if len(cust_id_to_subsession_map[customer_id]) > 0 {
			cust_id_to_subsession_map[customer_id] = append(cust_id_to_subsession_map[customer_id], output_strings[i])
		} else {
			cust_id_to_subsession_map[customer_id] = []string{output_strings[i]}
		}
	}
	for key, value := range cust_id_to_subsession_map {
		file_name := key + "-search-series-output.json"
		file, err := os.Create(file_name)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		_, err = file.Write([]byte("[" + strings.Join(value, "") + "]"))
		if err != nil {
			fmt.Println(err)
		}
	}
}
