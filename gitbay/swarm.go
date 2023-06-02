package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

	rows, err := ReadCSV("../cid_data.csv")
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
		peerCache = append(peerCache, row[0])
	}

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
			time.Sleep(time.Second * 1) // sleep so that we don"t spam the network with requests
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

	sh := shell.NewShell("localhost:5001")
	cids := []string{
		//	"QmUWAPdR2N7wzAhuMxpvhE69rge9ayKTnFB6WtFXiXoQMC", "QmauhUbPerCeRjJAbPCoDAdASuzWJ9sMrx5k7FUPeSQDyU", "QmR4kqBs7sis2dBuhfriAd4szr7VPERqnEvqvgaYvMomSj", "QmXLt7pG2XfnPXv98GiL2Y3xXJiaAkE8oA3yCT6L6n6J57", "QmVLMsn9dP6Uh2mCTdY2R5CDPKcZtnq7R6kcErmXohLDT2", "QmSyVtPo3GegNrDQj93gUKAMehghPTZV2oo6gfJd5yU9SQ", "QmYnfgAZfX111YhkJNN5JKUYv38T4s3NhZ1wJssAWCiExk", "QmYPytb5JpD6jQ7sPDdPRa9jQLnj8hjk4wCeDwZQDiNBgi", "Qmapmvi92Lf6yzHRj3uLuYXnPPL8QznkcBMUgijEtyWmyU", "Qmcubtd8GBz8ihk2t6nWTHmw9ChvABcsWgm5u4tt1XRxGC", "QmbyBDYob9rTHCxHgBhmr9t5FX7c8gMy2HnFz8EmRiCGrw", "Qma4TxpzdcYqDMzpgGAHYCyjrsbin2C6hTnXvjbMUFondJ", "QmYaMCtgVF46b5jTJ9n95F5yTgw9ZYCRNupBeJrvnphjTW", "QmaYukD44EuV5HeQBa6yyZbC15FWnMzhHrcaFfkuRpyD42", "QmepVMTaAzgx6rNdGA1kpvRDLtcJTq1Yxn6KmXg4RvRDp6",
		"QmS8DvEzaK37nCdkFAP5SFgJgcwF3oM6VRjsHVzhrToBat",
	}

	for i, cid := range cids {
		fmt.Println(i, len(cids))
		r, err := sh.Cat(cid + "/index.html")
		if err != nil {
			continue
		}
		data, err := ioutil.ReadAll(r)
		if err != nil {
			continue
		}
		cc.ToStorage(data, cid)
	}

}
