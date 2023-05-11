package main

import (
	"context"
	"fmt"

	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/bootstrap"
)

func main() {
	ctx := context.Background()
	conf := core.BuildCfg{}
	node, err := core.NewNode(ctx, &conf)
	if err != nil {
		panic(err)
	}
	peers, err := config.DefaultBootstrapPeers()
	if err != nil {
		panic(err)
	}
	bs_conf := bootstrap.BootstrapConfigWithPeers(peers)
	fmt.Println(bs_conf)
	err = node.Bootstrap(bs_conf) // this gets nil pointer dereference
	fmt.Println("THIS LINE WILL NOT PRINT")
	if err != nil {
		panic(err)
	}
}
