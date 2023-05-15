package main

import (
	"io/ioutil"
	"log"

	shell "github.com/ipfs/go-ipfs-api"
)

func dataFromCID(cid string) []byte {
	sh := shell.NewShell("localhost:5001")

	reader, err := sh.Cat(cid)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
