package main 

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    curl "github.com/andelf/go-curl"
    "github.com/gorilla/mux"

    _ "github.com/lib/pq"
    "database/sql"
)



type ConfigVars struct {
    DATABASE_URL    string  `json:"DATABASE_URL"`
    DevID       string      `json:"devID"`
    Mac1        string      `json:"mac1"`
    Mac2        string      `json:"mac2"`
    Mac3        string      `json:"mac3"`
    Mac4        string      `json:"mac4"`    
    Ticker_Key  string      `json:"ticker_key"`
}

var config_vars *ConfigVars
var IS_LIVE bool

func getConfigVarsLocal () bool {
    easy := curl.EasyInit()
    defer easy.Cleanup()

    respSuccess := func (buf []byte, userdata interface{}) bool {
        // println("DEBUG: size=>", len(buf))
        // println("DEBUG: content=>", string(buf))

        if err := json.Unmarshal(buf, config_vars); err != nil {
            panic(err)
        }
        return true
    }

    easy.Setopt(curl.OPT_URL, "https://api.heroku.com/apps/feelingcolor/config-vars")
    easy.Setopt(curl.OPT_NETRC, curl.NETRC_REQUIRED)
    easy.Setopt(curl.OPT_HTTPHEADER, []string{"Accept: application/vnd.heroku+json; version=3"})
    easy.Setopt(curl.OPT_WRITEFUNCTION, respSuccess)

    if err := easy.Perform(); err != nil {
        fmt.Printf("ERROR: %v\n", err)
        return false
    }

    return true
}


func getConfigVarsProd () bool {
    config_vars.DATABASE_URL = os.Getenv("DATABASE_URL")
    config_vars.DevID = os.Getenv("devID")
    config_vars.Mac1 = os.Getenv("mac1") 
    config_vars.Mac2 = os.Getenv("mac2")
    config_vars.Mac3 = os.Getenv("mac3")
    config_vars.Mac4 = os.Getenv("mac4")
    config_vars.Ticker_Key = os.Getenv("ticker_key")

    return true
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "9090"
    }

    router := mux.NewRouter()
    config_vars = &ConfigVars{}

    // get config vars
    db_url := os.Getenv("DATABASE_URL")
    if db_url == "" {
        getConfigVarsLocal()
    } else {
        getConfigVarsProd()
    }
    fmt.Printf("%+v\n", *config_vars)

    // new controller, log in
    _remoteLogin()

    // database
    var sql_err error
    db, sql_err = sql.Open("postgres", config_vars.DATABASE_URL)
    if sql_err != nil {
        log.Fatal(sql_err)
    }

    // API
    router.HandleFunc("/send-mood/{color}", SendMoodCommand).Methods("GET")
    router.HandleFunc("/bulb/color/{mac}/{color}", ChangeColorByMacCommand).Methods("GET")
    router.HandleFunc("/bulb/power/{mac}/{state}", SetPowerByMacCommand).Methods("GET")
    router.HandleFunc("/ctrl/start/{key}", StartTickerCommand).Methods("GET")
    router.HandleFunc("/ctrl/stop/{key}", StopTickerCommand).Methods("GET")

    // only expose if running localhost
    if os.Getenv("DATABASE_URL") == "" {
        // also ticker_key is just "209" if running locally
        config_vars.Ticker_Key = "209"

        router.HandleFunc("/ctrl/login", RemoteLoginCommand).Methods("GET")
        router.HandleFunc("/ctrl/get-devices", GetDevicesCommand).Methods("GET")
        router.HandleFunc("/ctrl/register/{mac}", RegisterDeviceCommand).Methods("GET")
        router.HandleFunc("/ctrl/deregister/{mac}", DeregisterDeviceCommand).Methods("GET")
    }

    // HTML
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
    http.Handle("/", router)

    // start off not live
    IS_LIVE = false

    fmt.Printf("Attempting to run server running on port " + port + "\n")
    err := http.ListenAndServe(":" + port, router) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }    
}
