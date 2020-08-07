# Get Bank Mega Promotions (Golang) #

This project is used to scrap [Bank Mega Promotions](https://www.bankmega.com/promolainnya.php) website.

## Getting Started

This project has been tested with [Go 1.14](https://golang.org/dl/) using [go.mod](https://blog.golang.org/using-go-modules).

## Installation
Install go package using:
```shell script
$ git clone https://github.com/fahimbagar/bankmega-promotions-go
 
$ go get -u
```

## Usage
- Run via go run
```shell script
$ go run ./...
```

- Or, build and run via binary executable
```shell script
$ go build -o solution

$  ./solution
```

## Unit Test
```shell script
$ go test ./... -v -count=1 
```

## Library
1. [htmlquery][htmlquery]
2. [assert][assert]

[htmlquery]: "https://github.com/antchfx/htmlquery"
[assert]: "github.com/stretchr/testify/assert"
[jest]: "https://github.com/facebook/jest"
