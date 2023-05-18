package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"

	cid "github.com/ipfs/go-cid"
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

func TestCidCollection() {
	dir_cids := [4]string{
		"bafybeicnxhkmocvutxrexcwj62eqidgunz22wqmwzrrghtdyvi5vjgn6ci",
		"QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o", // NOT A DIR, as a test case
		"QmQmhYjzuJUzsM3uMVtByzsfdQG6H3LeGTkUUD1yHVf1vb",
		"QmPXeM8QwpqBcnzE54fduPb5mm9trMDfd2adgL1KrmNNP6",
	}
	for _, c := range dir_cids {
		files, err := GetFileNamesFromCID(c)
		if err != nil {
			panic(err)
		}
		if len(files) != 0 {
			fmt.Println(c)
			fmt.Println(files)
			fmt.Println(GetDataFromCID(c + "/" + files[0]))
			fmt.Println("----")
		}
	}
}

func generateRandomCid() (cid.Cid, error) {
	length := 30
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return cid.Cid{}, err
	}
	_cid, err := cid.Builder.Sum(cid.V0Builder{}, bytes)
	return _cid, nil
}
