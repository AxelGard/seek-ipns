package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

type Crawler struct {
	host host.Host
	ctx  context.Context
	dht  *dht.IpfsDHT
}

const helloWorldCID = "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o"
const localKuboNodePeerID = "12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi"

func main() {
	ctx := context.Background()

	// Create a new libp2p host
	h, err := libp2p.New()
	if err != nil {
		log.Fatalf("Failed to create host: %v", err)
	}

	// Create a new DHT
	dht, err := dht.New(ctx, h)
	if err != nil {
		log.Fatalf("Failed to create DHT: %v", err)
	}

	// Bootstrap the DHT
	err = dht.Bootstrap(ctx)
	if err != nil {
		log.Fatalf("Failed to bootstrap DHT: %v", err)
	}
	cID, err := cid.Decode(helloWorldCID)
	if err != nil {
		panic(err)
	}
	peers, err := dht.FindProviders(ctx, cID)
	if err != nil {
		panic(err)
	}
	fmt.Println(peers)
}
