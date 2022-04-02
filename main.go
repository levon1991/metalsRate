package main

import (
	"fmt"

	"github.com/levon1991/metalsRate/dom"
)

func main() {
	d := dom.Init()
	fmt.Println(d.Pars("https://markets.businessinsider.com/commodities"))
}
