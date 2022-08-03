package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"strconv"
	"os"
	"encoding/json"
)

type pageConfig struct {
	Output		string
}

type appConfig struct {
	DaysOfWeek []int `json:"daysOfWeek"`
	SpecificDatesStr []string `json:"specificDates"`
	specificDates []time.Time
}

/* SETTINGS */

// serverPort indicates which tcp port to listen on and server http
const serverPort string = "8080"

// output time format
const timeFormat string = "02/01/2006"

// configPath contains the relative or absolute path to the event configuration being used to parse event details
const configPath string = "./configs/conf.json"

/* GLOBALS */ 

// pageConf saves the page's parsed config while making it available to the handlers
var pageConf pageConfig

// appConf saves the applications parsed config while making it available to other functions
var appConf appConfig

func calcEnd(start time.Time, duration int) (endDate time.Time) {
	var daysToSkip bool = true
	var oldEnd time.Time = start
	var interval int = duration
	var newEnd time.Time

	for daysToSkip {
		log.Print("new round")
		newEnd = oldEnd.AddDate(0, 0, interval)

		skips := 0
		for d := oldEnd; d.After(newEnd) == false; d = d.AddDate(0, 0, 1) {
			var alreadySkipped bool = false

			for _, dayOff := range appConf.DaysOfWeek {
				log.Printf("day off in conf %d", dayOff)
				log.Printf("day checking %d", int(d.Weekday()))
				if int(d.Weekday()) == dayOff{
					skips += 1
					alreadySkipped = true
					break
				}
			}

			if !alreadySkipped{
				for _, dateOff := range appConf.specificDates {
					if d == dateOff{
						skips += 1
						break
					}
				}
			}
		}

		if skips != 0{
			// shift the new interval checking for skip days to start after the current one (non-overlapping) as well as -1 from the skips to make the calc inclusive
			interval = skips - 1
			oldEnd = newEnd.AddDate(0, 0, 1)
		}else{
			daysToSkip = false
		}
	}

	return newEnd
}

func endHandler(w http.ResponseWriter, r *http.Request) {
	stringDate := fmt.Sprintf("%02s", r.FormValue("day")) + "/" + fmt.Sprintf("%02s", r.FormValue("month")) + "/" + r.FormValue("year")
	log.Printf(stringDate)

	startTime, err := time.Parse(timeFormat, stringDate)

	duration, err2 := strconv.Atoi(r.FormValue("duration"))

	if err != nil && err2 != nil{
		pageConf.Output = "Invalid date or length entered! Try again"
	}else {
		// the minus one on duration is to make it an inclusive of start/end calculation
		pageConf.Output = "End Date (dd/mm/yyyy): " + calcEnd(startTime, duration-1).Format(timeFormat)
	}

    t, _ := template.ParseFiles("./web/templates/index.html")
	t.Execute(w, pageConf) 
}

func homeHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("./web/templates/index.html")
	t.Execute(w, nil) 
}

func parseAppConf(){
	// Read the config file
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}	

	// Parse the configuration file and unmarshal into global config variable
	err = json.Unmarshal(configFile, &appConf)
	if err != nil {
		log.Fatal("Error during Unmarshal() of config ", err)
	}

	for _, day := range appConf.DaysOfWeek {
		if day >= int(time.Sunday) && day <= int(time.Saturday) {
			appConf.DaysOfWeek = append(appConf.DaysOfWeek, day)
		} else {
			log.Fatal("Invalid day of week found in config")
		}
	}

	for _, dateStr := range appConf.SpecificDatesStr {
		date, err := time.Parse(timeFormat, dateStr)

		if err != nil{
			log.Fatal("Invalid dd/mm/yyyy date found in config")
		}

		appConf.specificDates = append(appConf.specificDates, date)
	}
}

func main() {
	// Parse the configuration
	parseAppConf()

	// Serve the web application
	http.HandleFunc("/", homeHandler)
    http.HandleFunc("/end", endHandler)

	log.Printf("Starting server on 127.0.0.1:%s\n", serverPort)

	log.Fatalln(http.ListenAndServe(":" + serverPort, nil))
}