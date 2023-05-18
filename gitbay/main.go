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

	pID := "12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi"
	RecordKey := ipns.RecordKey(peer.ID(pID))
	sh := ipfs.GetIpfsShell()
	resolved, err := sh.Resolve(RecordKey)
	if err != nil {
		panic(err)
	}
	cid := strings.TrimPrefix(resolved, "/ipfs/")
	data := GetDataFromCID(cid)
	fmt.Println(data)

}

func GetDataFromIpnsUsingPeerID(peerID string, sh shell.Shell) string {
	RecordKey := ipns.RecordKey(peer.ID(peerID))
	resolved, err := sh.Resolve(RecordKey)
	if err != nil {
		panic(err)
	}
	cid := strings.TrimPrefix(resolved, "/ipfs/")
	return GetDataFromCID(cid)
}
