package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type TokenReponse struct {
	Access_Token             string `json:"access_token"`
	Token_Type               string `json:"token_type"`
	Expires_In               int    `json:"expires_in"`
	Refresh_Token            string `json:"refresh_token"`
	Refresh_Token_Expires_In int    `json:"refresh_token_expires_in"`
	Scope                    string `json:"scope"`
}

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

type DivisionLinkChannelResponse struct {
	overall_rank       int
	season_driver_data DivisionDriver
	link               string
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
	url := "https://members-ng.iracing.com/data/carclass/get"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", token_response.Token_Type+" "+token_response.Access_Token)
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
	fmt.Println(string(body))
}
