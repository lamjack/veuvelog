# Venvelog - Logging for Golang

## Installation

Install the latest version with

```bash
$ go get -u github.com/lamjack/veuvelog
```

## Basic Usage

```go
import (
	"github.com/lamjack/veuvelog"
	"github.com/lamjack/veuvelog/handler"
)

lvl := log.DEBUG
h := handler.NewConsoleHandler(logLevel)
ConsoleLogger = log.NewLogger("logger name")
ConsoleLogger.PushHandler(h)
```

## About

### Author

Jack Lam - <jack@wizmacau.com> - <http://lamjack.github.io/>

### License

Venvelog is licensed under the MIT License - see the `LICENSE` file for details