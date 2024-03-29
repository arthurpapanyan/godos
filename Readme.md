# HTTP Performance Testing Tool

This is a command-line tool written in Go for HTTP performance testing. It allows you to simulate HTTP requests to a specified URL with customizable concurrency, request count, request method, request body, and logging options.

## Installation

To install the tool, you need to have Go installed on your system. You can then install it using `go get`:

```bash
go get github.com/arthurpapanyan/godos
```


### Usage
```
Usage: ./godos [OPTIONS]

Options:
  -c int
        Concurrency level (default 1)
  -config string
        JSON config file path
  -d string
        Request body
  -logfile string
        Log file path
  -m string
        HTTP request method (default "GET")
  -n int
        Number of requests (default 1)
  -t string
        Target URL