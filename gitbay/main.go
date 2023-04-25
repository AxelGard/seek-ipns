package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
)

func main() {
	ctx := context.Background()

	// Create a libp2p host
	host, err := libp2p.New()
	if err != nil {
		fmt.Println("Failed to create libp2p host:", err)
		return
	}

	// Create a DHT instance
	ipfs_DHT, err := dht.New(ctx, host)
	if err != nil {
		fmt.Println("Failed to create DHT:", err)
		return
	}

	// Bootstrap the DHT to connect to the IPFS network
	if err = ipfs_DHT.Bootstrap(ctx); err != nil {
		fmt.Println("Failed to bootstrap DHT:", err)
		return
	}

	fmt.Println("DHT bootstrap successful")

	key, err := ipfs_DHT.GetPublicKey(ctx, host.ID())
	if err != nil {
		panic(err)
	}
	raw_key, err := key.Raw()
	str_key := crypto.ConfigEncodeKey(raw_key)

	peers, err := ipfs_DHT.GetClosestPeers(ctx, str_key)
	if err != nil {
		panic(err)
	}
	fmt.Println(peers)

}
