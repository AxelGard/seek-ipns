package main

import (
	"errors"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func RunCompare() {

	peers, err := LoadPeers(DATA_SET_OF_PEERS_TO_COMPARE_PATH)

	err = Compare(peers)
	if err != nil {
		log.Println("WARNING: Crashed trying to run compare: ", err)
		LogCrash()
	}
}

func LoadPeers(filePath string) ([]string, error) {
	rows, err := ReadCSV(filePath)
	if err != nil {
		return nil, err
	}
	peers := []string{}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		peers = append(peers, row[0])
	}
	return peers, nil
}

func Compare(peers []string) error {
	log.Println("Started running compearance of peers IPNS CIDs")

	sh := shell.NewShell("localhost:5001")
	ic := IpnsCollector{}
	ic.Init()
	timesIter := 0

	checkedPeers, err := LoadPeers(DATA_STORE_PATH_WEEK + WEEK_FILE)
	if err != nil {
		return err
	}

	if len(peers) == len(checkedPeers) {
		log.Println("checked peers is of equal length, will not run.\nCheck that new loaded file is correct.")
		return errors.New("No peers to check")
	}

	for len(checkedPeers) < len(peers) {
		log.Println("Have checked ", len(checkedPeers), "/", len(peers))
		for _, peer := range peers {
			if Contains(checkedPeers, peer) {
				continue
			}
			_, err := sh.FindPeer(peer)
			time.Sleep(time.Second * 5) // sleep so that we don"t spam the network with requests
			if err != nil {
				continue
			}
			log.Println("Have checked ", len(checkedPeers), "/", len(peers))
			cid, err := ic.GetIpnsCid(peer)

			row := []string{
				peer,
				time.Now().Format("2006-01-02 15:04:05"),
				cid,
			}
			err = AddRowToCSV(DATA_STORE_PATH_WEEK+WEEK_FILE, row)
			if err != nil {
				return nil
			}
			checkedPeers = append(checkedPeers, peer)
		}
		timesIter++
		if timesIter == TIMES_TRY_TO_CHECK_COMPARE {
			log.Println("Gave up trying to check all peers, all peers have not been checked")
			log.Println("Have checked ", len(checkedPeers), '/', len(peers))
			LogCrash("StoppedTrying")
			return nil
		}
		log.Println("Tried all peers will wait 30 min and try again.")
		time.Sleep(time.Minute * 30)
	}

	log.Println("ALL CHECKED, ALL DONE!")

	return nil
}
