## Details
This is a simple App to help me load test this other repos:
- Go API: https://github.com/Mhaxym/web-service-gin-docker
- NET CORE API: https://github.com/Mhaxym/core-cache-api

## How to use
- Clone this repository
- Run `go build` inside the main folder.
- Run `./api-client-load-testing {TEST_CODE}` where `TEST_CODE` is the code of the test we want to run. The accepted values are [`GO`, `NETCORE`].

For example, if you want to test the Go API, you should run `./api-client-load-testing GO`.

## Explanation
The code is pretty simple. It just creates a new goroutine for each request we want to make. Each goroutine will make a request to the API. The code is pretty self explanatory.

By default it will make 10000 requests to the API. You can change this value by changing the `CONCURRENT_REQUESTS` variable in the `main.go` file.

The output will also be generated in a log file.