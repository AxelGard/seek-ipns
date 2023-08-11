package main

import (
	ipns "github.com/ipfs/boxo/ipns"
	"github.com/libp2p/go-libp2p/core/peer"
)

func get_ipns_record_from_peer(peerID string) string {
	return ipns.RecordKey(peer.ID(peerID))
}
