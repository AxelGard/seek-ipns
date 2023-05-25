package main

import (
	"context"
	"fmt"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	//CrawlingEachNode()
	SwarmCrawl()
	//CheckGitProvders()
}

func SwarmCrawl() {
	ctx := context.Background()
	sh := shell.NewShell("localhost:5001")

	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	foundPeers := 0
	var peerCache []string

	for {
		info, err := sh.SwarmPeers(ctx)
		if err != nil {
			panic(err)
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
				//row := []string{peer, "False"}
				//AddRowToCSV("../peerlog.csv", row)
				continue
			}
			//row := []string{peer, "True"}
			//AddRowToCSV("../peerlog.csv", row)
			foundPeers++
			err = Collecting(peer, ic, cc)
			if err != nil {
				panic(err)
			}
		}

	}
}

func test_crawled_peer() {

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
	crawler.Init(bootstrapNodes, ctx, peerChan)
	go crawler.Run()
	sh := shell.NewShell("localhost:5001")
	peerNotFoundCount := 0
	FoundPeers := 0
	fmt.Println("Done with setup.")
	for peer := range peerChan {
		_, err := sh.FindPeer(peer)
		time.Sleep(time.Second * 2)
		if err != nil {
			peerNotFoundCount++
			continue
		}
		FoundPeers++
		fmt.Println("have ", FoundPeers, "/", FoundPeers+peerNotFoundCount)
	}
}
