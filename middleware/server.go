package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func startHttpServer() {
	http.HandleFunc("/getlatest", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		var count int
		var err error
		if len(params["count"]) != 0 && len(params["count"][0]) != 0 {
			count, err = strconv.Atoi(params["count"][0])
			if err != nil {
				http.Error(w, "Invalid count parameter", http.StatusBadRequest)
				return
			}
		} else {
			count = 1
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Write([]byte("["))
		for i := 0; i < count; i++ {
			if len(sensorBacklog)-1-i < 0 {
				break
			}
			if i != 0 {
				w.Write([]byte(","))
			}
			w.Write([]byte(sensorBacklog[len(sensorBacklog)-1-i]))
		}
		w.Write([]byte("]"))
		fmt.Printf("Sent %d entries to %s\n", count, r.RemoteAddr)
	})
	fmt.Println("Starting HTTP server on port 8080")
	http.ListenAndServe(":8080", nil)
}
