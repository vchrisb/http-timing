# http-timing
a golang application to get detailed http timings for an endpoint for troubleshooting purposes on Cloud Foundry or a different environment

## Cloud Foundry deployment

The application does issue a http request to the specified endpoint, adding a UUID as a query string, and will log long running requests.

### prepare

Copy `manifest.yml.example` to `manifest.yml` and adapt to your needs.

#### variables
`TEST_HTTP_ENDPOINT` The HTTP endpoint to be used for troubleshooting
`TEST_THRESHOLD` The threshold in milliseconds over which a log will written to stdout
`TEST_RATE` The rate of requests per second

#### dep

run `dep` to create vendoring files:
`dep ensure`

### push

push the application:
`cf push`

### evaluate logs

The application does write a log line with details, prepended with `CSV` when ever the threshold is reached.
One can filter the log with this tag and write to a file. The file can then imported with `space` as a delimiter.
`cf logs http-timing |  grep CSV  --line-buffered > app_log.csv`

Similarily on the target application, e.g. [cf-helloworld](https://github.com/vchrisb/cf-helloworld), you can filter on the router response time and use the UUID to correlate a slow request with the gorouter log.
`cf logs cf-helloworld | grep response_time --line-buffered > server_log.log`

## local deployment

```
dep ensure
go build
export TEST_HTTP_ENDPOINT="https://google.com"
export TEST_THRESHOLD=1000
export TEST_RATE=1
./http-timing
```
