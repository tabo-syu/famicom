package main

import "os"

const (
	success = 0
	failure = 1
)

func main() {
	if err := run(); err != nil {
		os.Exit(failure)
	}

	os.Exit(success)
}

func run() error {
	return nil
}
