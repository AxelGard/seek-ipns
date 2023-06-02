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
		"12D3KooWMXKtvVetid8sBR3AMaRNa2c5dTYabPeYWCr1SPuHFTeF",
		"12D3KooWHiV1Tq5dttWhgPyWkycJQkR2afrG66XRnpD7FDKHzrF5",
		"QmYGVgGGfD5N4Xcc78CcMJ99dKcH6K6myhd4Uenv5yJwiJ",
		"12D3KooWBJqN2X6K6qpT7zb1RNiVVoRKPG114W1wJShoriJ2hpZH",
		"12D3KooWCnnBKgQ7LhSUU2Fvw2qoJqicYxsmSPLTWexWZqeR7NTK",
		"QmVTUsvmcokhtW5sfXAqLm7Foh2n1zmQUunAhDkcAcyfTq",
		"12D3KooWHw6vFKpKdc82j2nsKJSSfhZQgpwYbwTKyWKMkZk9nRoN",
		"12D3KooWRLrYmJiotityz2JsvvivBqTdSFnNeh5NcKZwfikpGM7o",
		"12D3KooWGNWHDZpvTG3TzgYgtGsS365FzjiRah6eNhc98mqrnzNZ",
		"12D3KooWGDmMuvcceU8g3STFfFrh8JjDd3FgrMygGrfK5bWmRdY6",
		"12D3KooWCffjYJZ873aVipvNeXWBubQqgT3NQzZEg8EeQEefGUH5",
		"12D3KooWMbWYrXGiHy1NNYixDjQf5BCG3rMpgBMrUaP5By5sELvk",
		"12D3KooWCzpBAgnkmdCakWCpd9gfDwSUFqwGgVfFjJnoDAGHCgyC",
		"12D3KooWEo3jHX5KaZXdVPHW6doigj85jb4fRHgawZS8dVKRfJVx",
		"12D3KooWFrNCsnR9w8MfVcm2NEL3fAYBxVygAdfpAbtCqeKsJF8U",
		"12D3KooWBqgDYmbXdquvXUF46hVaKe4MVZLfYEutVRXarZqQQSZS",
		"12D3KooWADhw6mFtYVsFf1pZ4d8ZZKM71PBBAnbJtzXezkuGCTTC",
		"12D3KooWSPsLTZW9Gg7tHxEY2sNwXkzHqD9pqUhGWGrXGNXyLN8o",
		"12D3KooWJ2NmXCpKeAhku7scoyMAqjbF88DSLipUSrgovUGYo4wc",
		"12D3KooWLcck3DE33chSgEUdrSKB7UUWt3uayStL6gcVww7eDxSo",
		"12D3KooWLxDRpUNCY5SSH71MDRbUxMFwgZD7vXTV7cM6r71bcSwc",
		"12D3KooWDbeULzLoPgDt75v7jbWpX8h3VzCQUPX2stWoycYFCDGo",
		"12D3KooWQ15qxGi87cea4YKPc8bddJca8njaiU5WZcvprMwNaBfk",
		"12D3KooWAzNnK9YbvGjjBxxjyiXLQGhBFsuPzxpdVxBzAf2RUT4c",
		"12D3KooWKbe1Ukivkg8Z1HqxJeUWTTe2BrxzV53ur3egWXNZurPS",
		"12D3KooWJsk9pYxP7gxzGDG9uc3YbFAGniTYs6uxeziqzwN6wUwg",
		"12D3KooWSLyAoLYUrJMn4RXhSyAd63PhCnoF9Mhr7yqzVaABoYje",
		"12D3KooWAgh5H8TKTeXdpvQHJKySecrjnMGYCgw7j1wv7x2z5qd1",
		"12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi",
	}
	ctx := context.Background()
	sh := shell.NewShell("localhost:5001")
	for i, peer := range peers {
		fmt.Println(i, "/", len(peers))
		time.Sleep(time.Second * 1) // sleep so that we don"t spam the network with requests
		_, err := sh.FindPeer(peer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = Collecting(peer, ic, cc, &dc, ctx)
		fmt.Println(err)
	}

}
