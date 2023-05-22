package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ipfs/boxo/ipns"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p/core/peer"
)

type IpnsCollector struct {
	sh *shell.Shell
}

func (ic *IpnsCollector) Init() error {
	ic.sh = shell.NewShell("localhost:5001")
	return nil
}

func (ic *IpnsCollector) GetIpnsCid(peerId string) (string, error) {
	recordKey := ipns.RecordKey(peer.ID(peerId))
	resolved, err := ic.sh.Resolve(recordKey)
	if err != nil {
		return "", nil
	}
	cid := strings.TrimPrefix(resolved, "/ipfs/")
	return cid, nil
}

type CidCollector struct {
	sh *shell.Shell
}

func (cc *CidCollector) Init() error {
	cc.sh = shell.NewShell("localhost:5001")
	return nil
}

func (cc *CidCollector) GetDataFromCid(cid string) ([]byte, error) {
	reader, err := cc.sh.Cat(cid)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (cc *CidCollector) ToStorage(cid string) error {
	data, err := cc.GetDataFromCid(cid + "/README.md")
	if err != nil {
		return err
	}
	filePath := "../data/" + cid
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	f.Write(data)
	return nil
}

func (cc *CidCollector) GetFileNamesFromCid(cid string) ([]string, error) {
	var result []string
	f_list, err := cc.sh.List(cid)
	if err != nil {
		return nil, err
	}
	for _, f := range f_list {
		result = append(result, f.Name)
	}
	return result, nil
}

func test_GetFileNamesFromCid() {
	cc := CidCollector{}
	cc.Init()
	dir_cids := [4]string{
		"bafybeicnxhkmocvutxrexcwj62eqidgunz22wqmwzrrghtdyvi5vjgn6ci",
		"QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o", // NOT A DIR, as a test case
		"QmQmhYjzuJUzsM3uMVtByzsfdQG6H3LeGTkUUD1yHVf1vb",
		"QmPXeM8QwpqBcnzE54fduPb5mm9trMDfd2adgL1KrmNNP6",
	}
	for _, c := range dir_cids {
		files, err := cc.GetFileNamesFromCid(c)
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
