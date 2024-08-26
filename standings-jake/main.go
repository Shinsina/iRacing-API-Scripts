package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type MinimalInitialResponse struct {
	Link string `json:"link"`
}

type iRacingAPIResponse struct {
	Success           bool   `json:"success"`
	Season_ID         int    `json:"season_id"`
	Season_Name       string `json:"season_name"`
	Season_Short_Name string `json:"season_short_name"`
	Series_ID         int    `json:"series_id"`
	Series_Name       string `json:"series_name"`
	Car_Class_ID      int    `json:"car_class_id"`
	Race_Week_Num     int    `json:"race_week_num"`
	Division          int    `json:"division"`
	Club_ID           int    `json:"club_id"`
	Customer_Rank     int    `json:"customer_rank"`
	Chunk_Info        struct {
		Chunk_Size        int      `json:"chunk_size"`
		Num_Chunks        int      `json:"num_chunks"`
		Rows              int      `json:"rows"`
		Base_Download_Url string   `json:"base_download_url"`
		Chunk_File_Names  []string `json:"chunk_file_names"`
	} `json:"chunk_info"`
}

type ChunkResponseGrouping struct {
	Season_Name        string         `json:"season_name"`
	Overall_Rank       int            `json:"overall_rank"`
	Season_ID          int            `json:"season_id"`
	Car_Class_ID       int            `json:"car_class_id"`
	Division           int            `json:"division"`
	Season_Driver_Data DivisionDriver `json:"season_driver_data"`
	Division_Rank      int            `json:"division_rank"`
}

type DivisionDriver struct {
	Customer_ID  int    `json:"cust_id"`
	Division     int    `json:"division"`
	Rank         int    `json:"rank"`
	Display_Name string `json:"display_name"`
	Club_ID      int    `json:"club_id"`
	Club_Name    string `json:"club_name"`
	License      struct {
		Category_ID   int     `json:"category_id"`
		Category      string  `json:"category"`
		License_Level int     `json:"license_Level"`
		Safety_Rating float64 `json:"safety_rating"`
		IRating       int     `json:"irating"`
		Color         string  `json:"color"`
		Group_Name    string  `json:"group_name"`
		Group_ID      int     `json:"group_id"`
	} `json:"license"`
	Helmet struct {
		Pattern     int    `json:"pattern"`
		Color1      string `json:"color1"`
		Color2      string `json:"color2"`
		Color3      string `json:"color3"`
		Face_Type   int    `json:"face_type"`
		Helmet_Type int    `json:"helmet_type"`
	} `json:"helmet"`
	Weeks_Counted       int     `json:"weeks_counted"`
	Starts              int     `json:"starts"`
	Wins                int     `json:"wins"`
	Top5                int     `json:"top5"`
	Top25_Percent       int     `json:"top25_percent"`
	Poles               int     `json:"poles"`
	Avg_Start_Position  float64 `json:"avg_start_position"`
	Avg_Finish_Position float64 `json:"avg_finish_position"`
	Avg_Field_Size      float64 `json:"avg_field_size"`
	Laps                int     `json:"laps"`
	Laps_Led            int     `json:"laps_led"`
	Incidents           int     `json:"incidents"`
	Points              int     `json:"points"`
	Raw_Points          float64 `json:"raw_points"`
	Week_Dropped        bool    `json:"week_dropped"`
	Country_Code        string  `json:"country_code"`
	Country             string  `json:"country"`
}

type MinimalDivisionResponse struct {
	Customer_Rank int `json:"customer_rank"`
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
	raw_standings_input, err := os.ReadFile("./standings-input.json")
	if err != nil {
		fmt.Println(err)
	}
	var standings_input []string
	err = json.Unmarshal(raw_standings_input, &standings_input)
	if err != nil {
		fmt.Println(err)
	}
	channel_count := len(standings_input)
	channel := make(chan string, channel_count)
	for i := 0; i < channel_count; i++ {
		mapping := standings_input[i]
		mapping_split := strings.Split(mapping, "_")
		season_id, err := strconv.Atoi(mapping_split[0])
		if err != nil {
			fmt.Println(err)
		}
		car_class_id, err := strconv.Atoi(mapping_split[1])
		if err != nil {
			fmt.Println(err)
		}
		url := fmt.Sprintf("https://members-ng.iracing.com/data/stats/season_driver_standings?season_id=%d&car_class_id=%d", season_id, car_class_id)
		go func() {
			fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(standings_input)) + " standings")
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}
			for key, value := range cookies {
				req.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client := &http.Client{}
			// @todo Determine if this can be shortened even further
			sleep_time := 100 * i
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
			var minimal_initial_response MinimalInitialResponse
			err = json.Unmarshal(body, &minimal_initial_response)
			if err != nil {
				fmt.Println(err)
			}
			channel <- minimal_initial_response.Link
		}()
	}
	link_channel := make(chan iRacingAPIResponse, len(standings_input))
	for i := 0; i < len(standings_input); i++ {
		url := <-channel
		go func() {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}
			for key, value := range cookies {
				req.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client := &http.Client{}
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
				fmt.Println(err)
			}
			link_channel <- iracing_api_response
		}()
	}
	chunk_channel := make(chan ChunkResponseGrouping, channel_count)
	for i := 0; i < channel_count; i++ {
		channel_response := <-link_channel
		calculated_page_number := channel_response.Customer_Rank / channel_response.Chunk_Info.Chunk_Size
		page_number := calculated_page_number
		if (channel_response.Customer_Rank % channel_response.Chunk_Info.Chunk_Size) == 0 {
			page_number = page_number - 1
		}
		base_url := channel_response.Chunk_Info.Base_Download_Url
		url_extension := channel_response.Chunk_Info.Chunk_File_Names[page_number]
		url := base_url + url_extension
		go func() {
			fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(standings_input)) + " links")
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}
			for key, value := range cookies {
				req.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client := &http.Client{}
			resp, err := http_client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			var iracing_api_response []DivisionDriver
			err = json.Unmarshal(body, &iracing_api_response)
			if err != nil {
				fmt.Println(err)
			}
			me_index := slices.IndexFunc(iracing_api_response, func(driver DivisionDriver) bool { return driver.Customer_ID == 300752 })
			me := iracing_api_response[me_index]
			var chunk_channel_output ChunkResponseGrouping
			chunk_channel_output.Season_ID = channel_response.Season_ID
			chunk_channel_output.Overall_Rank = channel_response.Customer_Rank
			chunk_channel_output.Season_Name = channel_response.Season_Name
			chunk_channel_output.Car_Class_ID = channel_response.Car_Class_ID
			chunk_channel_output.Division = me.Division + 1
			chunk_channel_output.Season_Driver_Data = me
			chunk_channel <- chunk_channel_output
		}()
	}
	var output []ChunkResponseGrouping
	for i := 0; i < channel_count; i++ {
		chunk_channel_response := <-chunk_channel
		division_rank := 0
		if chunk_channel_response.Season_Driver_Data.Division > -1 {
			url := fmt.Sprintf("https://members-ng.iracing.com/data/stats/season_driver_standings?season_id=%d&car_class_id=%d&division=%d", chunk_channel_response.Season_ID, chunk_channel_response.Car_Class_ID, chunk_channel_response.Season_Driver_Data.Division)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}
			for key, value := range cookies {
				req.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client := &http.Client{}
			resp, err := http_client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			var minimal_initial_response MinimalInitialResponse
			err = json.Unmarshal(body, &minimal_initial_response)
			if err != nil {
				fmt.Println(err)
			}
			req2, err2 := http.NewRequest("GET", minimal_initial_response.Link, nil)
			if err2 != nil {
				fmt.Println(err2)
			}
			for key, value := range cookies {
				req2.AddCookie(&http.Cookie{Name: key, Value: value})
			}
			http_client2 := &http.Client{}
			resp2, err3 := http_client2.Do(req2)
			if err3 != nil {
				fmt.Println(err3)
			}
			defer resp.Body.Close()
			body2, err4 := io.ReadAll(resp2.Body)
			if err4 != nil {
				fmt.Println(err4)
			}
			var iracing_api_response MinimalDivisionResponse
			err = json.Unmarshal(body2, &iracing_api_response)
			if err != nil {
				fmt.Println(err)
			}
			division_rank = iracing_api_response.Customer_Rank
		}
		fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(standings_input)) + " chunks")
		var value ChunkResponseGrouping
		value.Season_Name = chunk_channel_response.Season_Name
		value.Overall_Rank = chunk_channel_response.Overall_Rank
		value.Season_ID = chunk_channel_response.Season_ID
		value.Car_Class_ID = chunk_channel_response.Car_Class_ID
		value.Division = chunk_channel_response.Division
		value.Season_Driver_Data = chunk_channel_response.Season_Driver_Data
		value.Division_Rank = division_rank
		output = append(output, value)
	}

	file_name := "jake-standings-output.json"
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file_output, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
	}
	_, err = file.Write(file_output)
	if err != nil {
		fmt.Println(err)
	}
}
