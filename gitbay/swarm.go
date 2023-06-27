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

	dht_crawler, _, err := SetupCrawler(ctx)

	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	dc := DataCollector{}
	dc.Init()

	foundPeers := 0
	var peerCache []string

	rows, err := ReadCSV("../store/cid_data.csv")
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
			go Collecting(peer, ic, cc, &dc, dht_crawler, ctx)

		}

	}
}

func recollectData() error {
	NO_FILES := []string{}
	ic := IpnsCollector{}
	ic.Init()

	cc := CidCollector{}
	cc.Init()

	dc := DataCollector{}
	dc.Init()

	//sh := shell.NewShell("localhost:5001")
	ctx := context.Background()
	cids := []string{
		"QmUVTKsrYJpaxUT7dr9FpKq6AoKHhEM7eG1ZHGL56haKLG",
		"QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn",
		"QmU14JgCGns98PQoZ3oENV559jGnsG3kXtCjEKoqTbM7j9",
		"Qmf4LMVhe94H3vBWnxBzDyaHoSUYnaRqXN5z18fUrwFQqJ",
		"QmUWAPdR2N7wzAhuMxpvhE69rge9ayKTnFB6WtFXiXoQMC",
		"QmWj827Jm1tvrz2gQmqtFjeALtvmZdfWHJfCwfYqX6QMvW",
		"QmRFzBRH43QcnyFL1T8ifYbFcobg11CueKUM4KKdbvWJEj",
		"bafybeibfjuheszjnsj4qwgfjn2y65r3h7gsmq4g3mllrdaq6wixnpd2paq",
		"QmQdgLZP3ahz4FcCR35VXZEeez1ZBK4xM5cCzAgJYSB1hp",
		"QmSpjaeFnZq3FEw8mbAj1Nj5nKkVmFSkhnV6R7TgJYzd9d",
		"QmUrTiPXCPcg7N3FuhEtqgHoS6JpfWnUMaY8xK82qsUve9",
		"bafybeiatmr27mek4uinkvh7ou7q4ryk7pdfq4toiuwqzufz4lwggxlwp74",
		"QmRR3kqk46rFU8zCpe63fLsPAtKpSo7TUHqL5cYv5voAjQ",
		"QmRBRRq8EfU38QjgxFQ1M7BDcn9Y87Tx8K19Jt2x6VEj3a",
		"bafybeigy3qiei4mpepfvq75hi5nzbwyjbg2viws5cjyzlkcb72jihoghs4",
		"QmTWaPf6394B3pZpBnGbHv4Jf79N3PjrK5v7bNXPU7oJxo",
		"QmaJegxhjiqho2fv6pTHCAnXjZxWg2jAy9q6m1zsR3tPsk",
		"QmNN37tKyeoR9iUt7E3XCMoVDB68gxhEicrZrA9FXir1sK",
		"bafkreie6kj5kgblmzzt7gtdqibiycwduwrklzrw54j5axyefe6bolcxy3a",
		"QmdzkuRhECvEmaMweSUEwRM2FA2XYw5ijBTkRp3oY2xncK",
	}

	peers := []string{
		"12D3KooWQvVfZMKUaMeDBS7QPHL2APY4ANDPUzyaCSesos4tgGxv",
		"12D3KooWCXYw5cw5NXwMwj6yuHwgvZ7wwP8NqENBBmEjee8ZFHi5",
		"Qmd2XA6Z7gw51SqcT8cMuBWDbuvrm63BaWZZE6ikS52Y7B",
		"12D3KooWSyKt3fiPKiEXoDJrqTBGpbab7xEB1qf4q9sUFzcy1gE4",
		"12D3KooWQtMJyD4mcP9Vi5Z96aMKvy9P6oxULFao7JuQmGaAEDgG",
		"12D3KooWHrv9S4EmoGSUDoBQgx3uSB9JbKq45tSsy7FqWgt2yAGL",
		"12D3KooWLGNuMB1ZJdrd2eFFQBsTsqQgLr8XBHNCSph59HXiqXnq",
		"12D3KooWB6Kt6tZ4uJGjD4hou7dVQawXgyKHrCTz74ZFEnwDd5Zi",
		"12D3KooWFeHu75n14FDrSUDXbb7cXhRUJCYxsRCQLA8nZ9R7uJEg",
		"12D3KooWGc7qCqwQvx9r96hwtmVhJSiXKK1qMFunXP3KiccJv64w",
		"12D3KooWEKm4yKmjYuwSihsZbw2Jpn4xP2Yw1UowZiucEuLMjNBm",
		"12D3KooWEhG6AASQxjC1SiUHzXn4HnE1GfJFhZ2eRxcvWvQZ6S47",
		"12D3KooWNdXJvLxG4zVshhJUDtU9Rc38PEAygz1bubvnPP4t78kr",
		"12D3KooWBu8HivPwrtYDQ8s8YTqm8LFouXGWf5jg1aQoimfo1noo",
		"12D3KooWDYpPdfCFf3CbKpcLNmyA2vmJs4JY55k8yje9R1MxSgdB",
		"12D3KooWFaioDnjkzfXUU8bUeNiXLKi5vpQsd5BZoRAKV9Zr5QJV",
		"12D3KooWL2FcDJ41U9SyLuvDmA5qGzyoaj2RoEHiJPpCvY8jvx9u",
		"12D3KooWKDQ7UkeL9TzTwLty6GuZzEkHdRm8qgcVLiiTZFV2EEfs",
		"12D3KooWBbkCD5MpJhMc1mfPAVGEyVkQnyxPKGS7AHwDqQM2JUsk",
		"12D3KooWMXKtvVetid8sBR3AMaRNa2c5dTYabPeYWCr1SPuHFTeF",
	}
	skip := "Qmf4LMVhe94H3vBWnxBzDyaHoSUYnaRqXN5z18fUrwFQqJ"
	for i, cid := range cids {
		if cid == skip {
			continue
		}

		peer := peers[i]
		fmt.Println(i, len(cids))

		files, err := cc.GetAllFilesNamesFromCid(cid, &dc, ctx)
		if err != nil {
			continue
		}
		if len(files) == 0 {
			cid_data, err := cc.GetDataFromCid(cid)
			if err != nil {
				continue
			}
			fmt.Println(cid, " --> ", string(cid_data))
			contentType, err := GetContentType(cid_data)
			if err != nil {
				cc.ToDiscovery(peer, cid, NO_FILES)
			}
			singel_file := []string{contentType}
			cc.ToDiscovery(peer, cid, singel_file)
			dc.GetFileData(cid, "NONE", ctx)
			dc.ToDiscovery()
			continue
		}
		cc.ToDiscovery(peer, cid, files)
		dc.ToDiscovery()
		fmt.Println(cid, " ==> ", files)
		if isGitRepo(files) {
			readme, err := GetRepoFileData(cid, *cc.sh, ctx)
			if readme == nil {
				log.Println("Found repository but failed to save, cid: ", cid)
				continue
			}
			err = cc.ToStorage(readme, cid)
			if err != nil {
				continue
			}
		}

	}
	return nil

}
