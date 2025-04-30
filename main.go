package main

import (
	"fmt"
	"log"

	"github.com/dtluna/nachmundtracker/business"
)

func main() {
	games, err := business.DecodeData()
	if err != nil {
		log.Fatalf("decoding data: %v\n", err)
	}

	fmt.Printf("%+v\n", games)
}
