package main

import (
	"fmt"
	"strings"

	"github.com/AxelGard/gitbay/ipfs"
	"github.com/ipfs/boxo/ipns"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {

	//Start_crawling()

	var _peers = [5]string{
		"12D3KooWBUj3rZqUTssrEoJv8nARq46C8e7V6tdYacdyhJ7baG5k",
		"12D3KooWJ9RrJTkmpT5ENQEUX4SDVSAFeTaxFQXqqsh2QqNKb9rq",
		"12D3KooWErNDthedEtwpQM7sCNcNVSjK5at71Ysu6NBr5a6wB2Sx",
		"12D3KooWJJ4J3jPvhGGLMKcMW1A1wgKiHuD6wfyHttVFEk5uN5vG",
		"12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi",
	}

	sh := ipfs.GetIpfsShell()
	for i, p := range _peers {
		_cid := GetDataFromIpnsUsingPeerID(p, *sh)
		fmt.Println(GetDataFromCID(_cid))
		fmt.Println(i, "------------")
	}
}

func GetDataFromIpnsUsingPeerID(peerID string, sh shell.Shell) string {
	RecordKey := ipns.RecordKey(peer.ID(peerID))
	resolved, err := sh.Resolve(RecordKey)
	if err != nil {
		fmt.Println("could not be resolved:", peerID)
		panic(err)
	}
	cid := strings.TrimPrefix(resolved, "/ipfs/")
	return cid
}
