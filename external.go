package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type robot struct {
	Model            string    `json:"model"`
	SerialNumber     string    `json:"serialNumber"`
	ManufacturedDate time.Time `json:"manufacturedDate"`
	Category         string    `json:"category"`
}

func hackIntoDbForRobotDetails() ([]robot, error) {
	resp, err := http.Get("https://robotstakeover20210903110417.azurewebsites.net/robotcpu")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var robots []robot
	err = json.Unmarshal(body, &robots)
	if err != nil {
		return nil, err
	}
	return robots, nil
}
