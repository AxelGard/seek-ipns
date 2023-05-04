package main

import (
	"context"
	"fmt"

	kubo "github.com/ipfs/kubo"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/bootstrap"
)

func main() {
	fmt.Println(kubo.GetUserAgentVersion())
	fmt.Println(config.DefaultConnMgrHighWater)

	ctx := context.Background()
	conf := core.BuildCfg{}
	node, err := core.NewNode(ctx, &conf)
	if err != nil {
		panic(err)
	}
	err = node.Bootstrap(bootstrap.DefaultBootstrapConfig)
	if err != nil {
		panic(err)
	}

}
