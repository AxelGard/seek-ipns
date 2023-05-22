package main

import (
	"context"
	"fmt"
)

func main() {
	test_run()
}

func test_run() {

	bootstrapNodes := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
		"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	} // NOTE: found nodes used ipfs deamon cmd $ipfs bootstrap list
	ctx := context.Background()
	crawler := Crawler{}
	err := crawler.Init(bootstrapNodes, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("DHT bootstrap successful")

	peers, err := crawler.GetClosestPeers()
	if err != nil {
		panic(err)
	}
	var peersIds []string
	for _, p := range peers {
		peersIds = append(peersIds, string(p.ID))
	}

	mypeer := "12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi"
	peersIds = append(peersIds, mypeer)

	ic := IpnsCollector{}
	ic.Init()
	var cids []string
	for _, p := range peersIds {
		c, err := ic.GetIpnsCid(p)
		if err != nil {
			panic(err)
		}
		if c != "" {
			cids = append(cids, c)
		}
	}
	fmt.Println(cids)
	cc := CidCollector{}
	cc.Init()
	for _, c := range cids {
		if c == "" {
			continue
		}
		files, err := cc.GetFileNamesFromCid(c)
		if err != nil {
			panic(err)
		}
		if len(files) == 0 {
			continue
		}
		if isGitRepo(files) {
			err = cc.ToStorage(c)
			if err != nil {
				panic(err)
			}
		}
	}

}
