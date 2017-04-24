package main 

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/dorothychen/light/api"
)

func main() {
    router := mux.NewRouter()
    
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
    http.Handle("/", router)

    fmt.Printf("Attempting to run server running on port 9090\n")
    err := http.ListenAndServe(":9090", router) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }    
}
