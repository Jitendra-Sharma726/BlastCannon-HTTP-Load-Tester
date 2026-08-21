package main

import (
  "net/http"
  "time"
)

//Result holds single request data
type Result struct {
  StatusCode int
  Latency    time.Duration
  Error      error
}

//Metrics holds aggregated test data
type Metrics struct {
  TotalTime   time.Duration
  AvgLatency  time.Duration
  MinLatency  time.Duration
  MaxLatency  time.Duration
  SuccessCount  int
  FailCount     int
  }

//RunLoadTest manages the worker pool and aggregates results
func RunLoadTest(url string, totalRequests int, concurrency int, onResult func(completed int)) Metrics {

  //Prevent zero-worker deadlocks
  if concurrency <= 0 {
    concurrency = 1
  }

  testStartTime := time.Now()

  //1. Configure HTTP client to prevent connection bottlenecks
  customClient := &http.Client{
    Transport: &http.Transport{
       MaxIdleConnsPerHost: concurrency,
      },
      Timeout: 10 * time.Second,
    }

  //2. Intialize buffered channels
  jobs := make(chan int, totalRequests)
  results := make(chan Result, totalRequests)

  //3.Spin up worker pool
  for i := 1; i <= concurrency; i++ {
    go worker(customClient, url, jobs, results)}
  }

//4.Dispatch jobs
for i := 1; i <= totalRequests; i++ {
  jobs <- i
}

//5.Signal workers to stop
close(jobs)

//6. Aggregate results
var totalLatency, minLatency, maxLatency time.Duration
var successCount, failCount int

for i := 1; i <= totalRequests; i++ {
  res := <-results

  //Trigger UI progress
  if onResult != nil {
    onResult(i)
  }
  
  //Count non-200 OK responses as failures
  if res.Error != nil || res.StatusCode != http.StatusOK {
    failCount++
    continue
  }

  //Tally successful requests
  successCount++
  totalLatency += res.Latency

  //Track max latency
  if res.Latency >   maxLatency {
    maxLatency = res.Latency
    }

  //Track min latency (0 handles the first success)
  if minLatency == 0 || res.Latency < minLatency {
    minLatency = res.Latency
    }
  }

//calculate average latency (guard against divide-by-zero)
var avgLatency time.Duration
if successCount > 0 {
  avgLatency = totalLatency / time.Duration(successCount)
}

return Metrics{
  TotalTime: time.Since(testStartTime),
  AvgLatency: avgLatency,
  MinLatency: minLatency,
  MaxLatency: maxLatency,
  SuccessCount: successCount,
  FailCount:    failCount,
  }
}

//worker pulls jobs from the channel and fires HTTP requests
func worker(client *http.Client, url string, jobs chan int, results chan Results) {

  //Process jobs until channel closes
  for range jobs {
    start := time.Now()

    resp, err := client.Get(url)
    latency := time.Since(start)

    //Handle network failures
    if err != nil {
      results <- Result{Error : err}
      continue
    }

    //Close body immediately (NO defer in loops)
    statusCode := resp.StatusCode
    resp.Body.Close()

    //Send success data
    results <- Result{
         StatusCode: statusCode,
         Latency: latency,
      }
    }
  }











            













  

  




  











  







    
