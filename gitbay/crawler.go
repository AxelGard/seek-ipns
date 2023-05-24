package main

import (
	"context"
	"fmt"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dht_crawler "github.com/libp2p/go-libp2p-kad-dht/crawler"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func CrawlingEachNode() {

	bootstrapNodes := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
		"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	} // NOTE: found nodes used ipfs deamon cmd $ipfs bootstrap list
	ctx := context.Background()

	peerChan := make(chan string)

	crawler := Crawler{}
	err := crawler.Init(bootstrapNodes, ctx, peerChan)
	if err != nil {
		panic(err)
	}
	fmt.Println("DHT bootstrap successful")

	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	var peerCache []string

	fmt.Println(Collecting("12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi", ic, cc))

	go crawler.Run()

	sh := shell.NewShell("localhost:5001")

	count := 0
	foundPeers := 0
	for peer := range peerChan {
		if Contains(peerCache, peer) {
			continue
		}
		count++
		_, err := sh.FindPeer(peer)
		time.Sleep(time.Second * 1) // sleep so that we don't spam the network with requests
		if err != nil {
			continue
		}
		foundPeers++
		log.Println("found peers", foundPeers, "/", count)
		Collecting(peer, ic, cc)
	}

}

func Collecting(peer string, ic IpnsCollector, cc CidCollector) error {
	cid, err := ic.GetIpnsCid(peer)
	if err != nil {
		return err
	}
	if cid == "" {
		return nil
	}
	files, err := cc.GetFileNamesFromCid(cid)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		cid_data, err := cc.GetDataFromCid(cid)
		if err != nil {
			return nil
		}
		fmt.Println(cid, " --> ", string(cid_data))
		cc.ToDiscovery(peer, cid, files)
		return nil
	}
	cc.ToDiscovery(peer, cid, files)
	fmt.Println(cid, " ==> ", files)
	if isGitRepo(files) {
		err = cc.ToStorage(cid)
		if err != nil {
			return err
		}
	}

	return nil

}

type Crawler struct {
	host              host.Host
	ipfs_DHT          *dht.IpfsDHT
	bootstrappedNodes []string
	discoverdPeers    []peer.AddrInfo
	ctx               context.Context
	peerChan          chan string
}

func (c *Crawler) Init(bootstrapNodes []string, ctx context.Context, peerChan chan string) error {
	var err error
	c.ctx = ctx
	c.peerChan = peerChan
	c.bootstrappedNodes = bootstrapNodes
	c.host, err = libp2p.New()
	if err != nil {
		return err
	}
	c.ipfs_DHT, err = dht.New(c.ctx, c.host)
	if err != nil {
		return err
	}

	err = c.ipfs_DHT.Bootstrap(c.ctx)
	if err != nil {
		return err
	}

	for _, peerStr := range c.bootstrappedNodes {
		p, err := multiaddr.NewMultiaddr(peerStr)
		if err != nil {
			return err
		}
		pInfo, err := peer.AddrInfoFromP2pAddr(p)
		if err != nil {
			return err
		}
		err = c.host.Connect(c.ctx, *pInfo)
		if err != nil {
			return err
		}
		err = c.ipfs_DHT.Bootstrap(c.ctx)
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *Crawler) GetClosestPeers() ([]*peer.AddrInfo, error) {
	key, err := c.ipfs_DHT.GetPublicKey(c.ctx, c.host.ID())
	if err != nil {
		return nil, err
	}
	raw_key, err := key.Raw()
	str_key := crypto.ConfigEncodeKey(raw_key)

	_peers, err := c.ipfs_DHT.GetClosestPeers(c.ctx, str_key)
	if err != nil {
		return nil, err
	}
	var peers_info []*peer.AddrInfo
	for _, p := range _peers {
		s, err := c.ipfs_DHT.FindPeer(c.ctx, p)
		if err != nil {
			continue
		}
		peers_info = append(peers_info, &s)
	}
	return peers_info, nil
}

func (c *Crawler) Run() error {
	ctx := context.Background()
	peers_info, err := c.GetClosestPeers()
	spider, err := dht_crawler.NewDefaultCrawler(c.host, dht_crawler.WithConnectTimeout(time.Second*10))
	if err != nil {
		return err
	}
	spider.Run(ctx, peers_info, dht_crawler.HandleQueryResult(func(p peer.ID, rtPeers []*peer.AddrInfo) { c.CrawlingResult(rtPeers) }), dht_crawler.HandleQueryFail(func(p peer.ID, err error) {}))
	return nil
}

func (c *Crawler) CrawlingResult(peers []*peer.AddrInfo) {
	for _, p := range peers {
		c.peerChan <- p.ID.String()
	}
}
