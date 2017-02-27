package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type prd struct {
    CurrentTime string `json:"tmstmp"`
    ArrivalTime string `json:"prdtm"`
    StopName    string `json:"stpnm"`
    StopID      string `json:"stpid"`
    Route       string `json:"rt"`
    VehicleID   string `json:"vid"`
    Message     string `json:"msg"`
}

type bustimeResponse struct {
    Data    map[string][]prd    `json:"bustime-response"`
}

func predictionsForStop(stopID string, apiKey string, results chan prd) {
    client := &http.Client{}

    req, err := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getpredictions", nil)
    if err != nil {
        log.Print(err)
        os.Exit(1)
    }

    queryParams := req.URL.Query()
    queryParams.Set("format", "json")
    queryParams.Set("key", apiKey)
    queryParams.Set("stpid", stopID)
    req.URL.RawQuery = queryParams.Encode()

    res, err := client.Do(req)
    if err != nil {
        log.Print(err)
        os.Exit(2)
    }
    defer res.Body.Close()

    br := &bustimeResponse{}
    json.NewDecoder(res.Body).Decode(br)

    for _, e := range br.Data["error"] {
        results <- e
    }
    for _, p := range br.Data["prd"] {
        results <- p
    }
}

func pollPredictions(stopIDs []string, apiKey string, results chan prd) {
    for {
        for _, stopID := range stopIDs {
            go predictionsForStop(stopID, apiKey, results)
        }
        time.Sleep(1 * time.Minute)
    }
}

func socketHandler(w http.ResponseWriter, r *http.Request)  {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    var apiKey = os.Getenv("API_KEY")
    var stopIDs = strings.Split(os.Getenv("STOP_IDS"), ",")
    var prdResults = make(chan prd)

    go pollPredictions(stopIDs, apiKey, prdResults)

    for p := range prdResults {
        err = conn.WriteJSON(p)
        if err != nil {
            log.Println(err)
            return
        }
    }
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", socketHandler)
    http.Handle("/", router)
    err := http.ListenAndServe("0.0.0.0:5000", nil)
    if err != nil {
        log.Fatal(err)
    }
}
