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
go build -o gambling-bot gambling-bot.go
```

### Usage & Configuration

```sh
cd cmd
./gambling-bot --help
```

Gambling Bot uses config files in .yaml format, see `config.yml` file inside
the `tests` directory for real life examples.

## Running the tests

```sh
go test github.com/Namarand/gambling-bot/internal/app
```

## Continuous Integration

See [drone.github.papey.fr/papey/gambling-bot/](https://drone.github.papey.fr/papey/gambling-bot/)

## Built With

- [go-twitch-irc](https://github.com/gempir/go-twitch-irc) - An IRC Twitch bot library

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

- **Namarand** - _Main author_ - [Namarand](https://github.com/Namarand)
- **Wilfried OLLIVIER** - _Main author_ - [Papey](https://github.com/papey)

## Licence

[LICENSE](LICENSE) file for details

## Notes

- An end user documentation is availailable at [gamble.docs.papey.fr](https://gamble.docs.papey.fr/)
- Real life bot usage example can be found at [www.twitch.tv/val_pl_magicarenafr](https://www.twitch.tv/val_pl_magicarenafr)
- [twitchapps.com/tmi/](https://twitchapps.com/tmi/) can be used to generate the OAuth token