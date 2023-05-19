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

	var _peers = [1]string{
		"12D3KooWMikyELrvaczTNDkWMGT2G3qrwagfLRfURADV6gAGWraQ",
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
