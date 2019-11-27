# Gambling Bot

A golang powered Twitch bot handling simple vote mechanism

## Getting Started

### Prerequisites

- [Golang](https://golang.org)

### Installing

Get deps using go mod

```sh
go mod vendor
```

then build

```sh
cd cmd
go build gambling-bot.go -o gambling-bot
```

## Running the tests

```sh
go test github.com/Namarand/gambling-bot/internal/app
```
