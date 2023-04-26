package main

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func main() {

	BootstrapNodes := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
		"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	} // NOTE: found nodes used ipfs deamon cmd $ipfs bootstrap list

	ctx := context.Background()

	// Create a libp2p host
	host, err := libp2p.New()
	if err != nil {
		fmt.Println("Failed to create libp2p host:", err)
		panic(err)
	}

	// Create a DHT instance
	ipfs_DHT, err := dht.New(ctx, host)
	if err != nil {
		fmt.Println("Failed to create DHT:", err)
		panic(err)
	}

	// Bootstrap the DHT to connect to the IPFS network
	if err = ipfs_DHT.Bootstrap(ctx); err != nil {
		fmt.Println("Failed to bootstrap DHT:", err)
		panic(err)
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
	fmt.Println("Closest Peers:")
	fmt.Println(peers)

	const c_id = "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o"
	cid, err := cid.Decode(c_id)
	prov_peers, err := ipfs_DHT.FindProviders(ctx, cid)
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Println("Nodes that have the set CID:")
	fmt.Println(prov_peers)

}
