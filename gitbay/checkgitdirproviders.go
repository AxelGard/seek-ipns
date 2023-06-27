package main

import (
	"context"
	"time"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
)

func CheckGitProvders() {
	bootstrapNodes := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
		"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	} // NOTE: found nodes used ipfs deamon cmd $ipfs bootstrap list
	ctx := context.Background()

	sh := shell.NewShell("localhost:5001")
	peerChan := make(chan string)
	crawler := Crawler{}
	crawler.Init(bootstrapNodes, ctx, peerChan)
	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	dc := DataCollector{}
	dc.Init()

	git_cid := "QmXm6JkdHVrpxHXkUahBHcRzpU4iZ5xAKJojeE5y3z1dMr"
	var git_cids []string
	git_sub_dir, err := sh.List(git_cid)
	if err != nil {
		panic(err)
	}
	for _, sub_dir := range git_sub_dir {
		git_cids = append(git_cids, sub_dir.Hash)
	}

	for _, git_sub_cid := range git_cids {
		CID, err := cid.Decode(git_sub_cid)
		if err != nil {
			panic(err)
		}
		peers, err := crawler.ipfs_DHT.FindProviders(ctx, CID)
		if err != nil {
			panic(err)
		}

		for _, p := range peers {
			peer := p.ID.String()
			_, err := sh.FindPeer(peer)
			time.Sleep(time.Second * 1) // sleep so that we don't spam the network with requests
			if err != nil {
				continue
			}
			err = Collecting(peer, ic, cc, &dc, &crawler, ctx)
			if err != nil {
				panic(err)
			}
		}

	}
}
