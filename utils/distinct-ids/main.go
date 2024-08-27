package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MinimalSubsession struct {
	Series_ID int `json:"series_id"`
	Season_ID int `json:"season_id"`
}

func main() {
	raw_subsessions_output, err := os.ReadFile("./subsessions-output.json")
	if err != nil {
		fmt.Println(err)
	}
	var minimal_subsessions []MinimalSubsession
	err = json.Unmarshal(raw_subsessions_output, &minimal_subsessions)
	if err != nil {
		fmt.Println(err)
	}
	distinct_series_ids_map := make(map[int]bool)
	distinct_season_ids_map := make(map[int]bool)
	for i := 0; i < len(minimal_subsessions); i++ {
		value := minimal_subsessions[i]
		distinct_series_ids_map[value.Series_ID] = true
		distinct_season_ids_map[value.Season_ID] = true
	}
	distinct_series_ids := make([]int, len(distinct_series_ids_map))
	distinct_season_ids := make([]int, len(distinct_season_ids_map))
	i := 0
	for key := range distinct_series_ids_map {
		distinct_series_ids[i] = key
		i++
	}
	j := 0
	for key := range distinct_season_ids_map {
		distinct_season_ids[j] = key
		j++
	}
	series_file_name := "distinct-series-ids-output.json"
	series_file, err := os.Create(series_file_name)
	if err != nil {
		fmt.Println(err)
	}
	defer series_file.Close()
	series_file_output, err := json.Marshal(distinct_series_ids)
	if err != nil {
		fmt.Println(err)
	}
	_, err = series_file.Write(series_file_output)
	if err != nil {
		fmt.Println(err)
	}
	season_file_name := "distinct-season-ids-output.json"
	season_file, err := os.Create(season_file_name)
	if err != nil {
		fmt.Println(err)
	}
	defer season_file.Close()
	file_output, err := json.Marshal(distinct_season_ids)
	if err != nil {
		fmt.Println(err)
	}
	_, err = season_file.Write(file_output)
	if err != nil {
		fmt.Println(err)
	}
}
