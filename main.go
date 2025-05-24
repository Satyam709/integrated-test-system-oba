package main

import (
    "fmt"
)
func main() {
    err := SetFakeTime(1747539000000) // Example timestamp in milliseconds
    if err != nil {
        fmt.Printf("Failed to set fake time: %v\n", err)
        return
    }
    fmt.Println("Successfully set fake time")
}
