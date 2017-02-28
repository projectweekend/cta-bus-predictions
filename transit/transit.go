package transit

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
)


type Service interface {
    predictions(output chan interface{})
}


type bustimeResponse struct {
    Data    map[string][]prd    `json:"bustime-response"`
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


type CTABusService struct {
    APIURL string
    APIKey string
    StopIDs []string
    Predictions chan interface{}
}


func (c *CTABusService) fetchForStop(stopID string) {
    client := &http.Client{}

    req, err := http.NewRequest("GET", c.APIURL, nil)
    if err != nil {
        log.Print(err)
        os.Exit(1)
    }

    queryParams := req.URL.Query()
    queryParams.Set("format", "json")
    queryParams.Set("key", c.APIKey)
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
        c.Predictions <- e
    }
    for _, p := range br.Data["prd"] {
        c.Predictions <- p
    }
}


func (c *CTABusService) FetchPredictions() {
    for {
        for _, stopID := range c.StopIDs {
            go c.fetchForStop(stopID)
        }
        time.Sleep(1 * time.Minute)
    }
}
