package main
 
import (
    "os"
    "fmt"
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/vikstrous/zengge-lightcontrol/control"
    "github.com/vikstrous/zengge-lightcontrol/remote"
)


var rc *remote.Controller
var controller *control.Controller

type Bulb struct {
    MAC     string      `json:"mac"`
    i       int         `json:"int"`
}

type ErrReponse struct {
    Err      string        `json:"err"`
}

type SuccessResponse struct {
    OK       bool           `json:"OK"`
}

var bulbs []Bulb

func GetDevicesCommand(w http.ResponseWriter, req *http.Request) {
    devices, err := rc.GetDevices()
    if err != nil {
        json.NewEncoder(w).Encode(err)
        return
    }
    json.NewEncoder(w).Encode(devices)
}

func RegisterDeviceCommand(w http.ResponseWriter, req *http.Request) {    
    params := mux.Vars(req)
    if params["mac"] == "" {
        err := "Please provide a mac address"
        json.NewEncoder(w).Encode(ErrReponse{Err: err})
        return
    }

    rc.RegisterDevice(params["mac"])
    json.NewEncoder(w).Encode(SuccessResponse{OK: true})
}

func DeregisterDeviceCommand(w http.ResponseWriter, req *http.Request) {    
    params := mux.Vars(req)
    if params["mac"] == "" {
        err := "Please provide a mac address"
        json.NewEncoder(w).Encode(ErrReponse{Err: err})
        return
    }

    rc.DeregisterDevice(params["mac"])
    json.NewEncoder(w).Encode(SuccessResponse{OK: true})
}

func ChangeColorByMacCommand(w http.ResponseWriter, req *http.Request) {    
    params := mux.Vars(req)
    err := ""
    if params["mac"] == "" {
        err = "Please provide a mac address"
    }
    if params["color"] == "" {
        err = "Please provide a color address"
    }
    if err != "" {
        json.NewEncoder(w).Encode(ErrReponse{Err: err})
        return
    }

    mac := params["mac"]
    color_str := params["color"]
    color := control.ParseColorString(color_str)
    if color == nil {
        json.NewEncoder(w).Encode(ErrReponse{Err: "Invalid color"})
        return
    }

    // set controller to this remote operator
    remote_ := remote.NewRemoteTransport(rc, mac)
    controller = &control.Controller{remote_}
    controller.SetColor(*color)

    resp := SuccessResponse{OK: true}
    json.NewEncoder(w).Encode(resp)
}

func SetPowerByMacCommand(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    err := ""
    if params["mac"] == "" {
        err = "Please provide a mac address"
    }
    if params["state"] != "ON" && params["state"] != "OFF" {
        err = "Please provide a valid state (ON/OFF)"
    }
    if err != "" {
        json.NewEncoder(w).Encode(ErrReponse{Err: err})
        return
    }

    // set controller to this remote operator
    remote_ := remote.NewRemoteTransport(rc, params["mac"] )
    controller = &control.Controller{remote_}

    if params["state"] == "ON" {
        controller.SetPower(true)
    }
    if params["state"] == "OFF" {
        controller.SetPower(false)
    }

    json.NewEncoder(w).Encode(SuccessResponse{OK: true})
}

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

    router.HandleFunc("/ctrl/get-devices", GetDevicesCommand).Methods("GET")
    router.HandleFunc("/ctrl/register/{mac}", RegisterDeviceCommand).Methods("GET")
    router.HandleFunc("/ctrl/deregister/{mac}", DeregisterDeviceCommand).Methods("GET")
    router.HandleFunc("/bulb/color/{mac}/{color}", ChangeColorByMacCommand).Methods("GET")
    router.HandleFunc("/bulb/power/{mac}/{state}", SetPowerByMacCommand).Methods("GET")

    fmt.Printf("Server running on port 8010\n")
    log.Fatal(http.ListenAndServe(":8010", router))
}


