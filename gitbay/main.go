package main

import (
	"context"
	"fmt"
	"log"
)

func main() {

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

	go crawler.Run()

	for peer := range peerChan {
		fmt.Println(peer)
		err := Collecting(peer, ic, cc)
		log.Println(err, peer)
	}

}

func Collecting(peer string, ic IpnsCollector, cc CidCollector) error {
	cid, err := ic.GetIpnsCid(peer)
	if err != nil {
		return err
	}
	if cid != "" {
		return nil
	}
	files, err := cc.GetFileNamesFromCid(cid)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	fmt.Println(cid, " : ", files)
	if isGitRepo(files) {
		err = cc.ToStorage(cid)
		if err != nil {
			return err
		}
	}

	return nil

}
