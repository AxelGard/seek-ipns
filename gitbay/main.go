package main

import (
	"context"
	"fmt"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	SwarmCrawl()
	//recollectData()
	//CrawlingEachNode()
	//CheckGitProvders()
	//test_GetFileNamesFromCid()
}

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
				continue
			}
			foundPeers++
			go Collecting(peer, ic, cc, &dc, ctx)

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

func recollectData() {
	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	dc := DataCollector{}
	dc.Init()

	peers := []string{
		"12D3KooWMXKtvVetid8sBR3AMaRNa2c5dTYabPeYWCr1SPuHFTeF",
		"12D3KooWBu8HivPwrtYDQ8s8YTqm8LFouXGWf5jg1aQoimfo1noo",
		"12D3KooWDYpPdfCFf3CbKpcLNmyA2vmJs4JY55k8yje9R1MxSgdB",
		"12D3KooWFeHu75n14FDrSUDXbb7cXhRUJCYxsRCQLA8nZ9R7uJEg",
		"12D3KooWL2FcDJ41U9SyLuvDmA5qGzyoaj2RoEHiJPpCvY8jvx9u",
		"12D3KooWHrv9S4EmoGSUDoBQgx3uSB9JbKq45tSsy7FqWgt2yAGL",
		"12D3KooWFaioDnjkzfXUU8bUeNiXLKi5vpQsd5BZoRAKV9Zr5QJV",
		"12D3KooWLGNuMB1ZJdrd2eFFQBsTsqQgLr8XBHNCSph59HXiqXnq",
		"12D3KooWEKm4yKmjYuwSihsZbw2Jpn4xP2Yw1UowZiucEuLMjNBm",
		"12D3KooWKDQ7UkeL9TzTwLty6GuZzEkHdRm8qgcVLiiTZFV2EEfs",
		"Qmd2XA6Z7gw51SqcT8cMuBWDbuvrm63BaWZZE6ikS52Y7B",
		"12D3KooWNdXJvLxG4zVshhJUDtU9Rc38PEAygz1bubvnPP4t78kr",
		"12D3KooWSyKt3fiPKiEXoDJrqTBGpbab7xEB1qf4q9sUFzcy1gE4",
		"12D3KooWCXYw5cw5NXwMwj6yuHwgvZ7wwP8NqENBBmEjee8ZFHi5",
		"12D3KooWQtMJyD4mcP9Vi5Z96aMKvy9P6oxULFao7JuQmGaAEDgG",
		"12D3KooWNWJoqPz23KMFKiXhoHKk8yUJVVSuFVE5wJBurcuLRQ6t",
		"12D3KooWL2FcDJ41U9SyLuvDmA5qGzyoaj2RoEHiJPpCvY8jvx9u",
		"12D3KooWG2qoPwip8t5LGa7fmRoRaX2TPBAxHiXTC5iJVmxij4oh",
		"12D3KooWPSxE9Uqi48jQyxToBoqqDCvDt36nd5FqCeLqWXwcm2LM",
		"12D3KooWCA5DZCctgiAVhoGfxbNVJJz3znRwJ5S3sc5aCbpE39et",
		"QmWaeLN18v1dQ5dC7iqGamaMHhrEbA9BhgNoRQTnuhWpfq",
		"12D3KooWFeHu75n14FDrSUDXbb7cXhRUJCYxsRCQLA8nZ9R7uJEg",
		"12D3KooWN5YkXpdcYvBhDpedP9YHdQDVzYuqDFgf2GJCqsVvgyTb",
		"12D3KooWPxCbs6bFm5Z5SAvmrgzSvEcumsvZJgaHiEf4ECz5xAhb",
		"12D3KooWP7r9RR2Z4XymZ3fzAYR45MZBw3Quep5eoVvxgc5gSVGQ",
		"12D3KooWMeY7CEitm4hdXiXrfdisGE94MAyRBwABF1edK8HXa9Kq",
		"12D3KooWFcBs3xNBrhEHYT6UVrWYgorkGWBeSsbSrJyMN7irUmHy",
		"12D3KooWAoXX38D4npL9DuxCpQ2TL3wbA5xPzFCz7L8LRJnAjATZ",
		"12D3KooWE2tgGhZyd3pDmgnjMmYtETczDT4vAhqxn8E1ezDq7Vu3",
		"12D3KooWAw2gTLa2LXhnx6tkhJcdPvvaYH82R6m8dxttDCPEkCtm",
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
