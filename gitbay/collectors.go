package main

import (
	"context"
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

func (cc *CidCollector) ToDiscovery(peer string, cid string, files []string) error {
	filepath := "../found.csv"
	filesAsStr := ""
	for _, f := range files {
		filesAsStr += f + ","
	}
	row := []string{
		peer,
		cid,
		filesAsStr,
	}
	err := AddRowToCSV(filepath, row)
	return err
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

func (cc *CidCollector) GetAllFilesNamesFromCid(cid string, ctx context.Context) ([]string, error) {
	var result []string
	f_list, err := cc.sh.List(cid)
	if err != nil {
		return nil, err
	}
	for _, f := range f_list {

		cidf := cid + "/" + f.Name
		if cidf[0] != '/' {
			cidf = "/ipfs/" + cidf
		}
		fs, err := cc.sh.FilesStat(ctx, cidf)
		if err != nil {
			return nil, err
		}
		if fs.Type == "directory" {
			sub_files, err := cc.GetAllFilesNamesFromCid(cidf, ctx)
			if err != nil {
				return nil, err
			}
			result = append(result, sub_files...)
		} else {
			result = append(result, f.Name)
		}
	}
	return result, nil
}

func test_GetFileNamesFromCid() {
	cc := CidCollector{}
	cc.Init()
	dir_cids := [5]string{
		"bafybeicnxhkmocvutxrexcwj62eqidgunz22wqmwzrrghtdyvi5vjgn6ci",
		"QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o", // NOT A DIR, as a test case
		"QmQmhYjzuJUzsM3uMVtByzsfdQG6H3LeGTkUUD1yHVf1vb",
		"QmPXeM8QwpqBcnzE54fduPb5mm9trMDfd2adgL1KrmNNP6",
		"QmNb2LcaN8hzSNp4g7z8FtLsqvNyo3XDiR1gnDna1TWMqe",
	}
	for _, c := range dir_cids {
		files, err := cc.GetFileNamesFromCid(c)
		if err != nil {
			panic(err)
		}
		if len(files) != 0 {
			fmt.Println(c)
			fmt.Println(files)
			fmt.Println("----")
		}
	}
}

func Collecting(peer string, ic IpnsCollector, cc CidCollector, ctx context.Context) error {
	EMPTY_CID := "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"
	NO_FILES := []string{}
	NO_CID := ""
	cid, err := ic.GetIpnsCid(peer)
	if err != nil {
		fmt.Println("ERROR::", err)
		return err
	}
	if cid == NO_CID {
		cc.ToDiscovery(peer, NO_CID, NO_FILES)
		return nil
	}
	if cid == EMPTY_CID {
		cc.ToDiscovery(peer, EMPTY_CID, NO_FILES)
	}
	files, err := cc.GetFileNamesFromCid(cid)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		cid_data, err := cc.GetDataFromCid(cid)
		if err != nil {
			return nil
		}
		fmt.Println(cid, " --> ", string(cid_data))
		contentType, err := GetContentType(cid_data)
		if err != nil {
			cc.ToDiscovery(peer, cid, NO_FILES)
		}
		if contentType == "text/plain; charset=utf-8," {
			contentType = "text/plain,"
		}
		singel_file := []string{contentType}
		cc.ToDiscovery(peer, cid, singel_file)
		return nil
	}
	cc.ToDiscovery(peer, cid, files)
	fmt.Println(cid, " ==> ", files)
	if isGitRepo(files) {
		err = cc.ToStorage(cid)
		if err != nil {
			return err
		}
	}

	return nil

}
