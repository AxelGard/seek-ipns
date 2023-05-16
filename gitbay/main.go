package main

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	peerId := "12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi"
	c := peer.ToCid(peer.ID(peerId)).String()
	fmt.Println(cid.Decode(c))
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
