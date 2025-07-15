package main

import (
	"fmt"

	"github.com/ehumba/blog-aggregator/internal/config"
)

func main() {
	// read config file
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	// set username
	err = cfg.SetUser("Sijie")
	if err != nil {
		panic(err)
	}
	// read again
	newCfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	//print output
	fmt.Println(newCfg)
}
