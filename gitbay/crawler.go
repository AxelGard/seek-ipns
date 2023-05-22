package main

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dht_crawler "github.com/libp2p/go-libp2p-kad-dht/crawler"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Crawler struct {
	host              host.Host
	ipfs_DHT          *dht.IpfsDHT
	bootstrappedNodes []string
	discoverdPeers    []peer.AddrInfo
	ctx               context.Context
}

func (c *Crawler) Init(bootstrapNodes []string, ctx context.Context) error {
	var err error
	c.ctx = ctx
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
		c.discoverdPeers = append(c.discoverdPeers, *p)
	}
}
