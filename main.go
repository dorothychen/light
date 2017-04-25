package main 

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    curl "github.com/andelf/go-curl"
    "github.com/gorilla/mux"
    "github.com/vikstrous/zengge-lightcontrol/remote"
)



type ConfigVars struct {
    DevID       string      `json:"devID"`
}

var config_vars *ConfigVars

func getConfigVars () bool {
    easy := curl.EasyInit()
    defer easy.Cleanup()

    respSuccess := func (buf []byte, userdata interface{}) bool {
        // println("DEBUG: size=>", len(buf))
        // println("DEBUG: content=>", string(buf))

        if err := json.Unmarshal(buf, config_vars); err != nil {
            panic(err)
        }
        fmt.Printf("%+v", *config_vars)
        return true
    }

    easy.Setopt(curl.OPT_URL, "https://api.heroku.com/apps/vast-crag-43585/config-vars")
    easy.Setopt(curl.OPT_NETRC, curl.NETRC_REQUIRED)
    easy.Setopt(curl.OPT_HTTPHEADER, []string{"Accept: application/vnd.heroku+json; version=3"})
    easy.Setopt(curl.OPT_WRITEFUNCTION, respSuccess)

    if err := easy.Perform(); err != nil {
        fmt.Printf("ERROR: %v\n", err)
        return false
    }

    return true
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "9090"
    }

    router := mux.NewRouter()
    config_vars = &ConfigVars{}

    // get device ID as command line parameter
    if len(os.Args) < 2 {
        // fmt.Println("Please specify a device ID")
        // return
        getConfigVars()
    } else {
        config_vars.DevID = os.Args[1]
    }
    
    // new controller, log in
    rc = remote.NewController("http://wifi.magichue.net/WebMagicHome/ZenggeCloud/ZJ002.ashx", "8ff3e30e071c9ef5b304d83239d0c707", config_vars.DevID)
    rc.Login()

    // API
    router.HandleFunc("/send-mood/{color}", SendMoodCommand).Methods("GET")
    router.HandleFunc("/ctrl/get-devices", GetDevicesCommand).Methods("GET")
    router.HandleFunc("/ctrl/register/{mac}", RegisterDeviceCommand).Methods("GET")
    router.HandleFunc("/ctrl/deregister/{mac}", DeregisterDeviceCommand).Methods("GET")
    router.HandleFunc("/bulb/color/{mac}/{color}", ChangeColorByMacCommand).Methods("GET")
    router.HandleFunc("/bulb/power/{mac}/{state}", SetPowerByMacCommand).Methods("GET")

    // HTML
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
    http.Handle("/", router)

    fmt.Printf("Attempting to run server running on port " + port + "\n")
    err := http.ListenAndServe(":" + port, router) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }    
}
