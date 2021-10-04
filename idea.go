package main

import "github.com/zaydek/retro/go/pkg/stdio_logger"

func main() {
	logger := stdio_logger.New(stdio_logger.LoggerOptions{Datetime: true})
	logger.Stdout("Hello, world!")
}
