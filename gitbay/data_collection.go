package main

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/pkg/errors"
)

type CidData struct {
	Cid      string
	Peer     string
	Size     uint64 // nr of Bytes
	Format   string
	FileName string
}

func (cd *CidData) Slice() []string {
	return []string{cd.Peer, cd.Cid, cd.FileName, strconv.FormatUint(cd.Size, 10), cd.Format}
}

func (cd *CidData) String() string {
	str := ""
	for _, d := range cd.Slice() {
		str += d + ", "
	}
	return str
}

func (cd *CidData) ToDiscovery() {
	filepath := DATA_STORE_PATH + "cid_data.csv"
	row := cd.Slice()
	err := AddRowToCSV(filepath, row)
	if err != nil {
		log.Println("failed to append to cid data, ", err)
	}
}

type DataCollector struct {
	sh           *shell.Shell
	CurrentFiles []CidData
	CurrentPeer  string
}

func (dc *DataCollector) Init() error {
	dc.sh = shell.NewShell("localhost:5001")
	return nil
}

func (dc *DataCollector) ToDiscovery() {
	for _, d := range dc.CurrentFiles {
		//fmt.Println("DEBUG - to discovery Data collector: ", d)
		d.ToDiscovery()
	}
}

func (dc *DataCollector) GetFileData(cid string, fileName string, ctx context.Context) (CidData, error) {
	data := CidData{}
	data.Cid = cid
	data.Peer = dc.CurrentPeer
	data.FileName = fileName
	if cid[0] != '/' {
		cid = "/ipfs/" + cid
	}
	fs, err := dc.sh.FilesStat(ctx, cid)
	if err != nil {
		return data, err
	}
	if fs.Type == "directory" {
		return data, errors.New("Data Collector GetFileData got CID for directory not file")
	}
	data.Size = fs.Size
	reader, err := dc.sh.Cat(cid)
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return data, err
	}
	format, err := GetContentType(d)
	if err != nil {
		return data, err
	}
	data.Format = format

	dc.CurrentFiles = append(dc.CurrentFiles, data)
	return data, nil
}
