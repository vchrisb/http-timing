package main

import ("github.com/tcnksm/go-httpstat"
        "time"
        "net/http"
        "github.com/asaskevich/govalidator"
        "io/ioutil"
        "io"
        "log"
        "os"
        "strconv"
        "fmt"
        "github.com/google/uuid"
      )


func main() {

  endpoint := os.Getenv("TEST_HTTP_ENDPOINT")
  if (!govalidator.IsURL(endpoint)) {
    log.Fatal("No or malformed URL specified!")
  }
  var limit time.Duration =  time.Millisecond * time.Duration(3000)
  test_limit, err := strconv.Atoi(os.Getenv("TEST_THRESHOLD"))
  if (err == nil) {
    limit =  time.Millisecond * time.Duration(test_limit)
  }

  var rate =  time.Second / time.Duration(1)
  test_rate, err := strconv.Atoi(os.Getenv("TEST_RATE"))
  if (err == nil) {
    rate =  time.Second / time.Duration(test_rate)
  }

  fmt.Println("HTTP Endpoint: ", endpoint)
  fmt.Println("Threshold: ", limit)
  fmt.Println("Rate: ", rate)

  throttle := time.Tick(rate)
  i := 1
  for i > 0 {
    <-throttle  // rate limit our Service.Method RPCs
    //fmt.Print(".")
    go call(endpoint, limit, i)
    i++
  }

}

func call(endpoint string, limit time.Duration, run int) {

  // generate uuid
  id := uuid.New()

  call_endpoint := endpoint + "?" + id.String()
  req, err := http.NewRequest("GET", call_endpoint , nil)
  if err != nil {
      log.Fatal(err)
  }


  // Create go-httpstat powered context and pass it to http.Request
  var result httpstat.Result
  ctx := httpstat.WithHTTPStat(req.Context(), &result)
  req = req.WithContext(ctx)

  tr := &http.Transport{
    DisableKeepAlives: true,
  }
  client := &http.Client {
    Transport: tr,
  }
  res, err := client.Do(req)
  if err != nil {
      log.Fatal(err)
  }

  if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
      log.Fatal(err)
  }
  res.Body.Close()

  end_time := time.Now()
  result.End(end_time)

  if (result.Total(end_time) > limit) {
    // Show results
    fmt.Println("### Total is ", result.Total(end_time))
    //fmt.Printf("%+v", result)
    fmt.Println("DNSLookup: ", result.DNSLookup.Seconds())
    fmt.Println("TCPConnection: ", result.TCPConnection.Seconds())
    fmt.Println("TLSHandshake: ", result.TLSHandshake.Seconds())
    fmt.Println("ServerProcessing: ", result.ServerProcessing.Seconds())
    fmt.Println("CSV", end_time.Format(time.RFC3339Nano), run, id, result.Total(end_time).Seconds(), result.DNSLookup.Seconds(), result.TCPConnection.Seconds(), result.TLSHandshake.Seconds(), result.ServerProcessing.Seconds())
  }

}
