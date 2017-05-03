package main
 
import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/vikstrous/zengge-lightcontrol/control"
    "github.com/vikstrous/zengge-lightcontrol/remote"

    _ "github.com/lib/pq"
    "database/sql"
    "fmt"
)


var rc *remote.Controller
var controller *control.Controller
var db *sql.DB

type ErrReponse struct {
    Err      string        `json:"err"`
}

type SuccessResponse struct {
    OK       bool           `json:"OK"`
}

type StatusResponse struct {
    IsLive  bool          `json:"is_live"`
    Light1  control.State        `json:"light1"`
    Light2  control.State        `json:"light2"`
    Light3  control.State        `json:"light3"`
    Light4  control.State        `json:"light4"`
}

func _checkBulbState () string {
    state, err := controller.GetState()
    ctrl_err := ""

    if err != nil {
        ctrl_err = err.Error()
    } else if !state.IsOn {
        ctrl_err = "Bulb is not on"
    }
    return ctrl_err
}

func RemoteLoginCommand(w http.ResponseWriter, req *http.Request) {
    _remoteLogin()
}

func _remoteLogin() {
    rc = remote.NewController("http://wifi.magichue.net/WebMagicHome/ZenggeCloud/ZJ002.ashx", "8ff3e30e071c9ef5b304d83239d0c707", config_vars.DevID)
    rc.Login()
}


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

    err_color := _changeColorByMac(mac, color_str)
    if err_color != nil {
        json.NewEncoder(w).Encode(err_color)
        return
    }

    resp := SuccessResponse{OK: true}
    json.NewEncoder(w).Encode(resp)
}

func _changeColorByMac(mac string, color_str string) (err *ErrReponse) {
    color := control.ParseColorString(color_str)
    if color == nil {
        err := &ErrReponse{Err: "Invalid color"}
        return err
    }

    // set controller to this remote operator
    remote_ := remote.NewRemoteTransport(rc, mac)
    controller = &control.Controller{remote_}
    
    ctrl_err := _checkBulbState()
    if ctrl_err != "" {
        err := &ErrReponse{Err: ctrl_err}
        return err
    }

    controller.SetColor(*color)  
    return nil
}

func _setPowerByMac(mac string, state string) (err *ErrReponse) {
    // set controller to this remote operator
    remote_ := remote.NewRemoteTransport(rc, mac)
    controller = &control.Controller{remote_}

    _, ctrl_err := controller.GetState()
    if ctrl_err != nil {
        return &ErrReponse{Err: ctrl_err.Error()}
    }

    if state == "ON" {
        controller.SetPower(true)
    }
    if state == "OFF" {
        controller.SetPower(false)
    }

    return nil
}

func SetPowerAllCommand(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    err := ""
    if params["state"] != "ON" && params["state"] != "OFF" {
        err = "Please provide a valid state (ON/OFF)"
    }
    if err != "" {
        json.NewEncoder(w).Encode(ErrReponse{Err: err})
        return
    }

    for i, mac := range [4]string{config_vars.Mac1, config_vars.Mac2, config_vars.Mac3, config_vars.Mac4} {
        err_power := _setPowerByMac(mac, params["state"])
        if err_power != nil {
            fmt.Printf("bulb %d: %s failed turning %s. %s\n", i+1, mac, params["state"], err_power)
        }        
    }

    resp := SuccessResponse{OK: true}
    json.NewEncoder(w).Encode(resp)

}

func SendMoodCommand(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    if params["color"] == "" {
        json.NewEncoder(w).Encode(ErrReponse{Err: "Invalid color"})
        return
    }
    
    // put color on queue
    var retColor string
    err := db.QueryRow(`INSERT INTO colors(timestamp, color) VALUES (CURRENT_TIMESTAMP, $1) RETURNING color;`, params["color"]).Scan(&retColor)
    if err != nil {
        // TODOOOO handle err
        json.NewEncoder(w).Encode(err)
        return
    }

    json.NewEncoder(w).Encode(SuccessResponse{OK: true})
}

func StartTickerCommand(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    if params["key"] != config_vars.Ticker_Key {
        json.NewEncoder(w).Encode(ErrReponse{Err: "Incorrect key"})
        return
    }

    startTicker()
    json.NewEncoder(w).Encode(SuccessResponse{OK: true})

}

func StopTickerCommand(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    if params["key"] != config_vars.Ticker_Key {
        json.NewEncoder(w).Encode(ErrReponse{Err: "Incorrect key"})
        return
    }

    stopTicker()
    json.NewEncoder(w).Encode(SuccessResponse{OK: true})
}

func StatusCommand(w http.ResponseWriter, req *http.Request) {
    resp := StatusResponse{IsLive: IS_LIVE}

    for i, mac := range [4]string{config_vars.Mac1, config_vars.Mac2, config_vars.Mac3, config_vars.Mac4} {
        remote_ := remote.NewRemoteTransport(rc, mac)
        controller = &control.Controller{remote_}
        state, _ := controller.GetState()

        if i == 0 {
            resp.Light1 = *state
        } else if i == 1{
            resp.Light2 = *state
        } else if i == 2 {
            resp.Light3 = *state
        } else if i == 3 {
            resp.Light4 = *state
        }
    }

    json.NewEncoder(w).Encode(resp)
}
