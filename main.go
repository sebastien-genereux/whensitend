package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type pageConfig struct {
	Output		string
}

/* SETTINGS */

// serverPort indicates which tcp port to listen on and server http
const serverPort string = "8080"

// output time format
const timeFormat string = "02/01/2006"

/* GLOBALS */ 

//configData saves the page's parsed config while making it available to the handlers
var configData pageConfig

func endHandler(w http.ResponseWriter, r *http.Request) {
	stringDate := fmt.Sprintf("%02s", r.FormValue("day")) + "/" + fmt.Sprintf("%02s", r.FormValue("month")) + "/" + r.FormValue("year")
	log.Printf(stringDate)

	startTime, err := time.Parse(timeFormat, stringDate)

	if err != nil {
		log.Printf(startTime.Format(timeFormat))
		configData.Output = "Invalid date entered! Try again"
	}else {
		configData.Output = "dd/mm/yyyy:" + startTime.Format(timeFormat)
	}

    t, _ := template.ParseFiles("./web/templates/index.html")
	t.Execute(w, configData) 
}

func main() {
	// Serve the web application
    http.HandleFunc("/end", endHandler)

	log.Printf("Starting server on 127.0.0.1:%s\n", serverPort)

	log.Fatalln(http.ListenAndServe(":" + serverPort, nil))
}