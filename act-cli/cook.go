package main

import (
	"log"
	"os"

	"github.com/thinxer/actnet"
)

func cook() {
	var (
		fin, fmodel *os.File
		err         error
	)
	fin, err = os.Open(*flagInput)
	check(err)
	defer fin.Close()

	m := actnet.NewModel()
	if _, err := os.Stat(*flagModel); err == nil {
		fmodel, err = os.OpenFile(*flagModel, os.O_RDWR, 0666)
		check(err)
		log.Println("loading...")
		check(m.Load(fmodel))
		check(fmodel.Close())
	} else if !os.IsNotExist(err) {
		check(err)
	}

	log.Println("cooking...")
	result, err := m.Merge(fin)
	log.Println(result, err)
	log.Println(m.String())

	log.Println("saving...")
	if fmodel, err = os.OpenFile(*flagModel, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666); err != nil {
		check(err)
	}
	check(m.Save(fmodel))
	check(fmodel.Close())
}
