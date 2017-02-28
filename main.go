package main

import (
    "log"
    "net/http"
    "os"
    "strings"
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    "github.com/projectweekend/cta-bus-predictions/transit"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func socketHandler(w http.ResponseWriter, r *http.Request)  {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    cta := &transit.CTABusService{
        APIURL: os.Getenv("API_URL"),
        APIKey: os.Getenv("API_KEY"),
        StopIDs: strings.Split(os.Getenv("STOP_IDS"), ","),
        Predictions: make(chan interface{}),
    }
    go cta.FetchPredictions()

    for p := range cta.Predictions {
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
