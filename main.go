package main

import (
	"fmt"
	"log"
	"somewherecosmic/aggregator/internal/config"
)

func main() {

	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	conf.SetUser()
	newConf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newConf)
}
