package main

import (
	"encoding/json"
	"errors"
	"github.com/gocolly/colly/v2"
	"io/ioutil"
	"leadzu/structs"
	"log"
	"os"
	"strings"
)

type Config structs.Config

func NewConfig() Config {
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

func main() {
	configuration := NewConfig()

	collector := colly.NewCollector()

	err := collector.Post("https://www.leadzu.com/login/check", map[string]string{"user": configuration.User , "pass": configuration.Password})
	if err != nil {
		log.Fatal(err)
	}

	collector.OnResponse(func(response *colly.Response) {
		if strings.Contains(string(response.Body), "logging_form") {
			log.Fatal("Error! User or password is not valid")
		}

		loginMode := !strings.Contains(response.Request.URL.String(), "updateGestionPresencia")
		if loginMode {
			log.Printf("Login success\n")
		} else {
			log.Printf("Success! You have imputed your work time today\n")
		}
	})

	collector.OnHTML("#hora_entrada", func(element *colly.HTMLElement) {
		_, err := func() (bool, error){
			if element.Attr("disabled")  == "disabled" {
				return false, errors.New("Error! You have already inputed today")
			}

			return  true, nil
		}()

		if err != nil {
			log.Fatal(err)
		}

		errPost := collector.Post("https://www.leadzu.com/user/updateGestionPresencia", map[string]string{
			"tipo_update": "today", "hora_entrada": configuration.StartHour, "hora_salida": configuration.EndHour,
			"horas_reales": configuration.RealHours ,"incidencias": ""})
		if errPost != nil {
			log.Fatal(errPost)
		}
	})

	collector.Visit("https://www.leadzu.com/user/profile#tabs-gestion-presencia")
}
