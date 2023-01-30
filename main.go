package main

import (
	"fmt"

	"github.com/MiyukiMori11/prettyAddressGen/internal/eth"
	"github.com/MiyukiMori11/prettyAddressGen/internal/generator"
	"github.com/MiyukiMori11/prettyAddressGen/internal/workerPool"
)

func main() {
	g, _ := generator.New(eth.New(), workerPool.New())

	fmt.Println(g.Generate(`[\w\d]+6666666$`, 3, 1))
}
