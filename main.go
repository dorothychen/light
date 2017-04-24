package main 

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/vikstrous/zengge-lightcontrol/remote"
)

func main() {
    router := mux.NewRouter()

    // get device ID as command line parameter
    if len(os.Args) < 2 {
        fmt.Println("Please specify a device ID")
        return
    }
    devID := os.Args[1]

    // new controller, log in
    rc = remote.NewController("http://wifi.magichue.net/WebMagicHome/ZenggeCloud/ZJ002.ashx", "8ff3e30e071c9ef5b304d83239d0c707", devID)
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

    fmt.Printf("Attempting to run server running on port 9090\n")
    err := http.ListenAndServe(":9090", router) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }    
}
