package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type MinimalInitialResponse struct {
	Link string `json:"link"`
}

type CSVChannelResponse struct {
	base_download_url string
	chunk_file_names  []string
	driver_index      int
	chunk_size        int
	season_name       string
	season_id         int
	car_class_id      int
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
	CSV_URL string `json:"csv_url"`
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

type DivisionLinkChannelReponse struct {
	overall_rank       int
	season_driver_data DivisionDriver
	link               string
}

type DivisionCSVChannelResponse struct {
	api_response       iRacingAPIResponse
	overall_rank       int
	season_driver_data DivisionDriver
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
	var input_file_path string
	fmt.Println("Enter the relative file path for the standings-input file")
	fmt.Scanln(&input_file_path)
	raw_standings_input, err := os.ReadFile(input_file_path)
	value := strings.Split(input_file_path, "~")
	cust_id := value[1]
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
	csv_channel := make(chan CSVChannelResponse, channel_count)
	for i := 0; i < channel_count; i++ {
		channel_response := <-link_channel
		csv_url := channel_response.CSV_URL
		go func() {
			fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(len(standings_input)) + " CSV links")
			req, err := http.NewRequest("GET", csv_url, nil)
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
			content, err := csv.NewReader(resp.Body).ReadAll()
			if err != nil {
				fmt.Println(err)
			}
			driver_index := slices.IndexFunc(content, func(driver_string []string) bool {
				// This matches on the custid column at the end of the row
				matched, err := regexp.MatchString(cust_id, driver_string[len(driver_string)-1])
				if err != nil {
					fmt.Println(err)
				}
				return matched
			})
			var csv_channel_output CSVChannelResponse
			csv_channel_output.chunk_size = channel_response.Chunk_Info.Chunk_Size
			csv_channel_output.driver_index = driver_index
			csv_channel_output.base_download_url = channel_response.Chunk_Info.Base_Download_Url
			csv_channel_output.chunk_file_names = channel_response.Chunk_Info.Chunk_File_Names
			csv_channel_output.season_id = channel_response.Season_ID
			csv_channel_output.season_name = channel_response.Season_Name
			csv_channel_output.car_class_id = channel_response.Car_Class_ID
			csv_channel <- csv_channel_output
		}()
	}
	chunk_channel := make(chan ChunkResponseGrouping, channel_count)
	for i := 0; i < channel_count; i++ {
		csv_channel_response := <-csv_channel
		calculated_page_number := csv_channel_response.driver_index / csv_channel_response.chunk_size
		page_number := calculated_page_number
		if (csv_channel_response.driver_index % csv_channel_response.chunk_size) == 0 {
			page_number = page_number - 1
		}
		base_url := csv_channel_response.base_download_url
		url_extension := csv_channel_response.chunk_file_names[page_number]
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
			cust_id_as_int, err := strconv.Atoi(cust_id)
			if err != nil {
				fmt.Println(err)
			}
			index := slices.IndexFunc(iracing_api_response, func(driver DivisionDriver) bool { return driver.Customer_ID == cust_id_as_int })
			if index == -1 {
				fmt.Println(csv_channel_response)
			}
			driver := iracing_api_response[index]
			var chunk_channel_output ChunkResponseGrouping
			chunk_channel_output.Season_ID = csv_channel_response.season_id
			chunk_channel_output.Overall_Rank = csv_channel_response.driver_index
			chunk_channel_output.Season_Name = csv_channel_response.season_name
			chunk_channel_output.Car_Class_ID = csv_channel_response.car_class_id
			chunk_channel_output.Division = driver.Division + 1
			chunk_channel_output.Season_Driver_Data = driver
			chunk_channel <- chunk_channel_output
		}()
	}
	var output []ChunkResponseGrouping
	division_channel_count := 0
	division_link_channel := make(chan DivisionLinkChannelReponse, channel_count)
	for i := 0; i < channel_count; i++ {
		chunk_channel_response := <-chunk_channel
		if chunk_channel_response.Season_Driver_Data.Division > -1 {
			url := fmt.Sprintf("https://members-ng.iracing.com/data/stats/season_driver_standings?season_id=%d&car_class_id=%d&division=%d", chunk_channel_response.Season_ID, chunk_channel_response.Car_Class_ID, chunk_channel_response.Season_Driver_Data.Division)
			division_channel_count += 1
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
				var division_link_channel_output DivisionLinkChannelReponse
				division_link_channel_output.overall_rank = chunk_channel_response.Overall_Rank
				division_link_channel_output.season_driver_data = chunk_channel_response.Season_Driver_Data
				division_link_channel_output.link = minimal_initial_response.Link
				division_link_channel <- division_link_channel_output
			}()
		} else {
			var value ChunkResponseGrouping
			value.Season_Name = chunk_channel_response.Season_Name
			value.Overall_Rank = chunk_channel_response.Overall_Rank
			value.Season_ID = chunk_channel_response.Season_ID
			value.Car_Class_ID = chunk_channel_response.Car_Class_ID
			value.Division = chunk_channel_response.Division
			value.Season_Driver_Data = chunk_channel_response.Season_Driver_Data
			value.Division_Rank = 0
			output = append(output, value)
		}
	}
	division_csv_channel := make(chan DivisionCSVChannelResponse, division_channel_count)
	for i := 0; i < division_channel_count; i++ {
		go func() {
			division_link_channel_response := <-division_link_channel
			req, err := http.NewRequest("GET", division_link_channel_response.link, nil)
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
			var division_csv_channel_output DivisionCSVChannelResponse
			division_csv_channel_output.api_response = iracing_api_response
			division_csv_channel_output.overall_rank = division_link_channel_response.overall_rank
			division_csv_channel_output.season_driver_data = division_link_channel_response.season_driver_data
			division_csv_channel <- division_csv_channel_output
		}()
	}

	for i := 0; i < division_channel_count; i++ {
		division_rank := 0
		division_csv_channel_response := <-division_csv_channel
		fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(division_channel_count) + " CSV links")
		req, err := http.NewRequest("GET", division_csv_channel_response.api_response.CSV_URL, nil)
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
		content, err := csv.NewReader(resp.Body).ReadAll()
		if err != nil {
			fmt.Println(err)
		}
		driver_index := slices.IndexFunc(content, func(driver_string []string) bool {
			// This matches on the custid column at the end of the row
			matched, err := regexp.MatchString(cust_id, driver_string[len(driver_string)-1])
			if err != nil {
				fmt.Println(err)
			}
			return matched
		})
		division_rank = driver_index
		fmt.Println(strconv.Itoa(i) + " of " + strconv.Itoa(division_channel_count) + " chunks")
		var value ChunkResponseGrouping
		value.Season_Name = division_csv_channel_response.api_response.Season_Name
		value.Overall_Rank = division_csv_channel_response.overall_rank
		value.Season_ID = division_csv_channel_response.api_response.Season_ID
		value.Car_Class_ID = division_csv_channel_response.api_response.Car_Class_ID
		value.Division = division_csv_channel_response.api_response.Division
		value.Season_Driver_Data = division_csv_channel_response.season_driver_data
		value.Division_Rank = division_rank
		output = append(output, value)
	}

	file_name := fmt.Sprintf("~%s~standings-output.json", cust_id)
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
