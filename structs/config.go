package structs

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	StartHour string `json:"startHour"`
	EndHour   string `json:"endHour"`
	RealHours string `json:"realHours"`
}

func NewConfig() Config {
	configByInput, err := GetConfigByFlagInput(); if err != nil {
		log.Fatal(err)
	}

	if configByInput != (Config{}) {
		return configByInput
	}

	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully Opened config.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	return config
}

func GetConfigByFlagInput() (Config, error) {
	inputMode := false
	inputs := make(map[string]string)

	inputs["user"] = *flag.String("u", "", "user string")
	inputs["password"] = *flag.String("p", "", "password string")
	inputs["startHour"] = *flag.String("sh", "", "start hours string")
	inputs["endHour"] = *flag.String("eh", "", "end hours string")
	inputs["realHours"] = *flag.String("rh", "", "real hours string")

	for _, val := range inputs {
		if val != "" {
			inputMode = true
			break
		}
	}

	if !inputMode {
		return Config{}, nil
	}

	if inputs["user"] == "" || inputs["password"] == "" {
		return Config{}, errors.New("user and password must be provided")
	}

	config := NewConfig()
	config.User = inputs["user"]
	config.Password = inputs["password"]

	if inputs["startHour"] != "" {
		config.StartHour = inputs["startHour"]
	}

	if inputs["endHour"] != "" {
		config.EndHour = inputs["endHour"]
	}

	if inputs["realHours"] != "" {
		config.RealHours = inputs["realHours"]
	}

	return config, nil
}