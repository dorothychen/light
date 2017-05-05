package main 

import (
    "fmt"
    "log"
    "time"
)

var lastTimestamp string
var Q []string
var lights [4]string
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

func initColors() [4]string {
    rows, err := db.Query("SELECT color FROM colors ORDER BY timestamp LIMIT 4;")
    if err != nil {
        fmt.Println(err)
    }
    defer rows.Close()

    var colors [4]string
    i := 0
    for rows.Next() {
        err := rows.Scan(&colors[i])
        i += 1

        if err != nil {
            // TODOOOOOO probs just continue tbh
        }
    }

    for ; i < 4; i++ {
        colors[i] = "000000"
    }

    // TODO what does this actually accomplish; look this up 
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    return colors
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
    if IS_LIVE {
        return
    }

    lastTimestamp = time.Now().UTC().Format("2006-01-02 15:04:0000")
    fmt.Println("starting time:", lastTimestamp)

    // init light colors
    lights = initColors()

    // init ticker 
    ticker = time.NewTicker(2 * time.Second)

    IS_LIVE = true

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

func stopTicker() {
    if !IS_LIVE {
        return
    }
    ticker.Stop()
    IS_LIVE = false
}