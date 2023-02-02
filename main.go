package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/MiyukiMori11/prettyAddressGen/internal/eth"
	"github.com/MiyukiMori11/prettyAddressGen/internal/generator"
	"github.com/MiyukiMori11/prettyAddressGen/internal/workerPool"
)

var network = flag.String("network", "eth", "defines network")
var workers = flag.Int64("workers", 1, "changes workers count")
var results = flag.Int64("results", 1, "changes results count")
var pattern = flag.String("pattern", `[\w\d]+$`, "changes address pattern")

func main() {
	flag.Parse()

	var nw generator.AddrCreator

	//TODO: add validator
	switch *network {
	case "eth":
		nw = eth.New()
	default:
		log.Fatal("unexpected network")
	}

	g, err := generator.New(nw, workerPool.New())

	if err != nil {
		panic(err)
	}

	r, err := g.Generate(*pattern, *workers, *results)
	if err != nil {
		panic(err)
	}

	for _, a := range r {
		fmt.Printf("address: %s\npublic key: %s\nprivate key: %s\n", a.Address(), a.PublicKey(), a.PrivateKey())
	}
}
