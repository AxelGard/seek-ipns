package main

import (
	"context"
	"fmt"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func SwarmCrawl() {
	ctx := context.Background()
	sh := shell.NewShell("localhost:5001")

	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	dc := DataCollector{}
	dc.Init()

	foundPeers := 0
	var peerCache []string

	for {
		info, err := sh.SwarmPeers(ctx)
		if err != nil {
			log.Println("Swarm Peers failed!! will wait for 5 min and try again.")
			time.Sleep(time.Minute * 5)
			continue
		}
		fmt.Println("-------")
		log.Println("Starting to crawl new swarm ")
		fmt.Println("Found new peers: ", foundPeers)
		fmt.Println("Cache size: ", len(peerCache))
		fmt.Println(".........")
		for _, swarmNode := range info.Peers {
			peer := swarmNode.Peer
			if Contains(peerCache, peer) {
				continue
			}
			peerCache = append(peerCache, peer)

			_, err := sh.FindPeer(peer)
			time.Sleep(time.Second * 1) // sleep so that we don't spam the network with requests
			if err != nil {
				continue
			}
			foundPeers++
			go Collecting(peer, ic, cc, &dc, ctx)

		}

	}
}

func recollectData() {
	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	dc := DataCollector{}
	dc.Init()

	peers := []string{
		"12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi",
	}
	ctx := context.Background()
	sh := shell.NewShell("localhost:5001")
	for i, peer := range peers {
		fmt.Println(i, "/", len(peers))
		if i <= 5 {
			continue
		}
		time.Sleep(time.Second * 1) // sleep so that we don't spam the network with requests
		_, err := sh.FindPeer(peer)
		if err != nil {
			continue
		}
		err = Collecting(peer, ic, cc, &dc, ctx)
		fmt.Println(err)
	}

}
