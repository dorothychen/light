package main 

import (
    "fmt"
    "log"
    "time"
)

var lastTimestamp string
var Q []string
var lights []string
var ticker *time.Ticker

func getColorsFromDb() {
    rows, err := db.Query("SELECT * FROM colors WHERE timestamp > $1;", lastTimestamp)
    if err != nil {
        // TODOOOOOO
        fmt.Println(err)
    }
    defer rows.Close()

    var t string
    var color string
    for rows.Next() {
        err := rows.Scan(&t, &color)
        Q = append(Q, color)
        if err != nil {
            // TODOOOOOO probs just continue tbh
        }
        
        // track last timestamp seen
        lastTimestamp = t
    }

    // TODO what does this actually accomplish; look this up 
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
}

func updateColors() string {
    if len(Q) > 0 {
        new_col := Q[0]
        Q = Q[1:]

        lights[3] = lights[2]
        lights[2] = lights[1]
        lights[1] = lights[0]
        lights[0] = new_col

        // send these colors to the lights
        _changeColorByMac(config_vars.Mac1, lights[0])
        _changeColorByMac(config_vars.Mac2, lights[1])
        _changeColorByMac(config_vars.Mac3, lights[2])
        _changeColorByMac(config_vars.Mac4, lights[3])

        return new_col
    }
    return ""
}

func startTicker() {
    // init light colors
    lights = []string{"000000", "000000", "000000", "000000"}

    // init ticker 
    ticker = time.NewTicker(2 * time.Second)
    lastTimestamp = time.Now().UTC().Format("2006-01-02 15:04:0000")
    fmt.Println("starting time:", lastTimestamp)

    go func() {
        for t := range ticker.C {
            _ = t
            getColorsFromDb()
            updateColors()

            fmt.Printf("%v\n", lights)
            fmt.Println("ticked:", t, "\t lastTimestamp: " + lastTimestamp)
        }
    }()
}

