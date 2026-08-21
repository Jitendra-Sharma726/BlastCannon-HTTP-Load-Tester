package main

import (
  "fmt"
  "net/http"
  "time"
  )

const (
  TargetURL    = "http://localhost:8080"
  TotalRequests = 200
  Concurrency   = 20
)

func main() {
  fmt.Println("=== HTTP Load Tester ===")
  fmt.Printf("Target:
