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
  fmt.Printf("Target: %s\n", TargetURL)
  fmt.Printf("Request: %d | Concurrent Workers: %d\n", TotalRequests, Concurrency)
  fmt.Println("----------------------------------------------------------------------")

  //Spin up mock server in the background
  go startMockServer()

  //Wait for server boot
  time.Sleep(1 * time.Second)

  fmt.Println("Firing requests...")
  fmt.Print("Progress: ")

  //Print a dot every 10% of total requests
  interval := TotalRequests / 10

  metrics := RunLoadTest(TargetURL, TargetRequests, Concurrency, func(completed int) {
    if completed%interval == 0 {
      fmt.Print(".")
      }
    })

  fmt.Println()

  //Calculate RPS and Error Rate
  rps := float64(TotalRequests) / metrics.TotalTime.Seconds()
  errorRate := float64(metrics.FailCount) / float64(TotalRequests) * 100

  //Print final report
  fmt.Println("\n-------------------------------------------------------")
  fmt.Println("=== Load Test Results ===")
  fmt.Println("-------------------------------------------------------")
  fmt.Printf("Total Time:               %.2f seconds\n", metrics.TotalTime.Seconds())
  fmt.Printf("Request/Second:           %.2f RPS\n", rps)
  fmt.Println("--------------------------------------------------------")
  fmt.Println("Avg Latency:             %v\n", metrics.AvgLatency)
  fmt.Println("Min Latency:             %v\n", metrics.MinLatency)
  fmt.Println("Max Latency:             %v\n", metrics.MaxLatency)
  fmt.Println("--------------------------------------------------------")
  fmt.Println("Success (200s):         %d\n", metrics.SuccessCount)
  fmt.Printf("Failures:                %d\n", metrics.FailCount)
  fmt.Printf("Error Rate:              %.1f%%\n", errorRate)
  fmt.Println("----------------------------------------------------")
}

//HELPER FUNCTION : Mock Local Server
func startMockServer() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

    //Simulate a 15ms response time
    time.Sleep(15 * timeMillisecond)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
  })
  _ = http.ListenAndServe(":8080", nil)
  }












             















  
