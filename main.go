package main

import (
    "log"
    Routers "es/routers"
)



func main() {
    router := Routers.SetupRouter()

    if err := router.Run(":8083"); err != nil {
        log.Fatal(err)
    }
}