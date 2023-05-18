package main

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dht_crawler "github.com/libp2p/go-libp2p-kad-dht/crawler"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func Start_crawling() ([]*peer.AddrInfo, error) {

	BootstrapNodes := [6]string{
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
		return nil, err
	}

	// Create a DHT instance
	ipfs_DHT, err := dht.New(ctx, host)
	if err != nil {
		fmt.Println("Failed to create DHT:", err)
		return nil, err
	}

	// Bootstrap the DHT to connect to the IPFS network
	if err = ipfs_DHT.Bootstrap(ctx); err != nil {
		fmt.Println("Failed to bootstrap DHT:", err)
		return nil, err
	}

	for _, peerStr := range BootstrapNodes {
		p, err := multiaddr.NewMultiaddr(peerStr)
		if err != nil {
			return nil, err
		}
		pInfo, err := peer.AddrInfoFromP2pAddr(p)
		if err != nil {
			return nil, err
		}
		err = host.Connect(ctx, *pInfo)
		if err != nil {
			return nil, err
		}
		err = ipfs_DHT.Bootstrap(ctx)
		if err != nil {
			return nil, err
		}

	}

	fmt.Println("DHT bootstrap successful")

	key, err := ipfs_DHT.GetPublicKey(ctx, host.ID())
	if err != nil {
		return nil, err
	}
	raw_key, err := key.Raw()
	str_key := crypto.ConfigEncodeKey(raw_key)

	_peers, err := ipfs_DHT.GetClosestPeers(ctx, str_key)
	if err != nil {
		return nil, err
	}
	fmt.Println("Closest Peers:")
	fmt.Println(_peers)

	fmt.Println("")

	spider, err := dht_crawler.NewDefaultCrawler(host, dht_crawler.WithConnectTimeout(time.Second*10))
	if err != nil {
		fmt.Println("error init crawler ")
		return nil, err
	}
	var peers_info []*peer.AddrInfo
	for i, p := range _peers {
		s, err := ipfs_DHT.FindPeer(ctx, p)
		if err != nil {
			fmt.Println("error formating peer adr, idx: ", i)
		} else {

			peers_info = append(peers_info, &s)
		}
	}
	spider.Run(ctx, peers_info, dht_crawler.HandleQueryResult(func(p peer.ID, rtPeers []*peer.AddrInfo) { crawlQueryResult(p) }), dht_crawler.HandleQueryFail(func(p peer.ID, err error) {}))

	//const c_id = "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o"
	//cid, err := cid.Decode(c_id)
	//prov_peers, err := ipfs_DHT.FindProviders(ctx, cid)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println("")
	//fmt.Println("Nodes that have the set CID:")
	//fmt.Println(prov_peers)
	return peers_info, nil

}

func crawlQueryResult(p peer.ID) {
	cid := peer.ToCid(p)
	fmt.Println(cid)
}
