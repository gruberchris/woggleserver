# Woggle Server

Woggle is a word toggle TCP server written in Go. The Woggle Server accepts a word and returns it in reverse.

## Usage

```shell
$ go run main.go
```

## Connecting

```shell
$ nc localhost 3000
````

## Docker

```shell
$ docker build -t woggleserver .
$ docker run -p 3000:3000 woggleserver
```
