package utils

import (
	"encoding/json"
	"log"
	"strconv"
)

// DecodeDropDown returns ids of value selections. Use DecodeLabels to list Index value for each dropdown label.
func DecodeDropDown(valueIn string) []string {
	var val struct {
		IDs []int `json:"ids"`
	}
	err := json.Unmarshal([]byte(valueIn), &val)
	if err != nil {
		log.Println("DecodeDropDown Unmarshal Failed, ", err)
		return nil
	}
	result := make([]string, len(val.IDs))
	for i, id := range val.IDs {
		result[i] = strconv.Itoa(id)
	}
	return result
}

// DecodeLabels displays index value of all labels for a column. Uses column settings_str (see GetColumns).
// Use for Status(color) and Dropdown fields.
func DecodeLabels(settingsStr, columnType string) {
	var statusLabels struct {
		Labels         map[string]string `json:"labels"`             // index: label
		LabelPositions map[string]int    `json:"label_positions_v1"` // index: position
	}
	type dropdownEntry struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var dropdownLabels struct {
		Labels []dropdownEntry `json:"labels"`
	}

	if columnType == "color" {
		err := json.Unmarshal([]byte(settingsStr), &statusLabels)
		if err != nil {
			log.Fatal("DecodeLabels Failed", err)
		}
	}
	if columnType == "dropdown" {
		err := json.Unmarshal([]byte(settingsStr), &dropdownLabels)
		if err != nil {
			log.Fatal("DecodeLabels Failed", err)
		}
	}
}

func BuildDate(date string) DateTime {
	return DateTime{Date: date}
}
func BuildDateTime(date, time string) DateTime {
	return DateTime{Date: date, Time: time}
}
func BuildStatusIndex(index int) StatusIndex {
	return StatusIndex{index}
}
func BuildCheckbox(checked string) Checkbox {
	return Checkbox{checked}
}
