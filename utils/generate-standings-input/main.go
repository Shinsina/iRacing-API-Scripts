package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MinimalSubsession struct {
	Car_Class_ID int `json:"car_class_id"`
	Season_ID    int `json:"season_id"`
}

func main() {
	var input_file_path string
	fmt.Println("Enter the relative file path for the subsessions-output file")
	fmt.Scanln(&input_file_path)
	raw_subsessions_output, err := os.ReadFile(input_file_path)
	if err != nil {
		fmt.Println(err)
	}
	value := strings.Split(input_file_path, "-")
	cust_id := value[0]
	var minimal_subsessions [][]MinimalSubsession
	err = json.Unmarshal(raw_subsessions_output, &minimal_subsessions)
	if err != nil {
		fmt.Println(err)
	}
	distinct_season_car_class_ids_map := make(map[string]bool)
	for i := 0; i < len(minimal_subsessions); i++ {
		for j := 0; j < len(minimal_subsessions[i]); j++ {
			value := minimal_subsessions[i][j]
			concatenated_value := strconv.Itoa(value.Season_ID) + "_" + strconv.Itoa(value.Car_Class_ID)
			distinct_season_car_class_ids_map[concatenated_value] = true
		}
	}
	distinct_season_car_class_ids := make([]string, len(distinct_season_car_class_ids_map))
	i := 0
	for key := range distinct_season_car_class_ids_map {
		distinct_season_car_class_ids[i] = key
		i++
	}
	file_name := fmt.Sprintf("~%s~standings-input.json", cust_id)
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file_output, err := json.Marshal(distinct_season_car_class_ids)
	if err != nil {
		fmt.Println(err)
	}
	_, err = file.Write(file_output)
	if err != nil {
		fmt.Println(err)
	}
}
