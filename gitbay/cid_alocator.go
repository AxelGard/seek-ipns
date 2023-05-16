package main

import (
	"io/ioutil"
	"log"

	shell "github.com/ipfs/go-ipfs-api"
)

func GetDataFromCID(cid string) string {
	sh := shell.NewShell("localhost:5001")
	reader, err := sh.Cat(cid)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

func GetFileNamesFromCID(cid string) ([]string, error) {
	var result []string
	sh := shell.NewShell("localhost:5001")
	f_list, err := sh.List(cid)
	if err != nil {
		return nil, err
	}
	for _, f := range f_list {
		result = append(result, f.Name)
	}
	return result, nil
}
