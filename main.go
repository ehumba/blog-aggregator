package main

import (
	"github.com/ehumba/blog-aggregator/internal/config"
)

func main() {
	// read config file
	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	currentState := state{
		cfg: &conf,
	}

}
