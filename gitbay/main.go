package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func main() {

	BootstrapNodes := []string{
		"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/dnsaddr/bootstrap.libp2p.io/ipfs/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/ipfs/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	} // found nodes in a NPM pkg https://www.npmjs.com/package/libp2p-bootstrap

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

	for _, peerStr := range BootstrapNodes {
		p, err := multiaddr.NewMultiaddr(peerStr)
		if err != nil {
			panic(err)
		}
		pInfo, err := peer.AddrInfoFromP2pAddr(p)
		if err != nil {
			panic(err)
		}
		err = host.Connect(ctx, *pInfo)
		if err != nil {
			panic(err)
		}
		err = ipfs_DHT.Bootstrap(ctx)
		if err != nil {
			panic(err)
		}

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
