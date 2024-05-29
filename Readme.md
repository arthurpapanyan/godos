# HTTP Performance Testing Tool

This is a command-line tool written in Go for HTTP performance testing. It allows you to simulate HTTP requests to a specified URL with customizable concurrency, request count, request method, request body, and logging options.

## Installation

To install the tool, you need to have Go installed on your system. You can then install it using `go get`:

```bash
#Make sure that go bin directory included in your shell config

export PATH="$HOME/go/bin:$PATH"
```

```bash
go install github.com/arthurpapanyan/godos@latest
```


### Usage
```
Usage: godos [OPTIONS]

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
```

### Examples

```sh
#Sending 10 GET request with 2 parallel group. Total 20 requests
godos -t https://example.com -m GET -n 10 -c 2 

# Sending 5 POST request with
godos -c 1 -n 5 -d '{"email":"example@gmail.com","password":"Qwer1234!"}' -t http://example.com -m POST

```