package main

import (
	"os"

	"github.com/thinxer/actnet"
)

func show() {
	fin, err := os.Open(*flagModel)
	check(err)

	m := actnet.NewModel()
	check(m.Load(fin))
	m.PrintStats()
}
