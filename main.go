package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

type Data struct {
	Status Status `json:"status"`
}

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func init() {
	go autoReloadJSON()
}

func main() {
	http.HandleFunc("/", autoReloadWeb)
	http.ListenAndServe(":8080", nil)
}

func autoReloadJSON() {
	for {
		file, _ := os.Open("status.json")
		defer file.Close()

		data := Data{
			Status: Status{},
		}
		dbyte, err := ioutil.ReadAll(file)
		if err = json.Unmarshal(dbyte, &data); err != nil {
			fmt.Println("error unmarshal: ", err)
			return
		}

		data.Status.Water = rand.Intn(20)
		data.Status.Wind = rand.Intn(20)

		// fmt.Println(data.Status.Water)
		// fmt.Println(data.Status.Wind)
		newJson, _ := json.Marshal(data)

		if err = ioutil.WriteFile("status.json", newJson, 0644); err != nil {
			fmt.Println("error writefile: ", err)
			return
		}
		file.Sync()

		time.Sleep(time.Second * 15)
	}
}

func autoReloadWeb(w http.ResponseWriter, r *http.Request) {

	// read file JSON
	file, _ := os.Open("status.json")
	defer file.Close()

	data := Data{
		Status: Status{},
	}
	dbyte, err := ioutil.ReadAll(file)
	if err = json.Unmarshal(dbyte, &data); err != nil {
		fmt.Println("error unmarshal: ", err)
		return
	}

	var statusWater, statusWind string
	if data.Status.Water <= 5 {
		statusWater = "aman"
	} else if data.Status.Water <= 8 {
		statusWater = "siaga"
	} else {
		statusWater = "bahaya"
	}

	if data.Status.Wind <= 6 {
		statusWind = "aman"
	} else if data.Status.Wind <= 15 {
		statusWind = "siaga"
	} else {
		statusWind = "bahaya"
	}

	// fmt.Println(string(dbyte))
	fmt.Println("statusWater", statusWater)
	fmt.Println("statusWind", statusWind)
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Println("error http request: ", err)
		return
	}

	var finalData = map[string]string{
		"wind":        strconv.Itoa(data.Status.Wind),
		"statusWind":  statusWind,
		"water":       strconv.Itoa(data.Status.Water),
		"statusWater": statusWater,
	}

	t.Execute(w, finalData)
}
