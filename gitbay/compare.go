package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func RunCompare() {

	peers, err := LoadPeers(DATA_SET_OF_PEERS_TO_COMPARE_PATH_AGAINST)

	err = Compare(peers)
	if err != nil {
		log.Println("WARNING: Crashed trying to run compare: ", err)
		LogCrash()
	}
}

func LoadPeers(filePath string) ([]string, error) {
	const EMPTY_CID = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"
	rows, err := ReadCSV(filePath)
	if err != nil {
		return nil, err
	}
	peers := []string{}
	empty_peers := []string{}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		cid := row[2]
		if cid != EMPTY_CID && cid != "NONE" {
			peers = append(peers, row[0])
		} else {
			empty_peers = append(empty_peers, row[0])
		}
	}
	all_peers := append(peers, empty_peers...)
	return all_peers, nil
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
		logIterations(timesIter)
		for _, peer := range peers {
			if Contains(checkedPeers, peer) {
				continue
			}
			_, err := sh.FindPeer(peer)
			time.Sleep(time.Second * 1) // sleep so that we don"t spam the network with requests
			if err != nil {
				continue
			}
			log.Println("Have checked ", len(checkedPeers), "/", len(peers))
			go ArchivePeer(peer, ic)
			checkedPeers = append(checkedPeers, peer)
		}
		checkedAllOnes()
		timesIter++
		log.Println("ALL peers have been checked", timesIter, "times")
		if STOP_CHECK_AFTER_ITER && timesIter == TIMES_TRY_TO_CHECK_COMPARE {
			log.Println("Gave up trying to check all peers, all peers have not been checked")
			log.Println("Have checked ", len(checkedPeers), '/', len(peers))
			LogCrash("StoppedTrying")
			return nil
		}
		log.Println("Tried all peers will wait 1 hour and automatically try again.")
		time.Sleep(time.Hour * 1)
	}

	log.Println("ALL CHECKED, ALL DONE!")
	checkedAllOnes()
	return nil
}

func ArchivePeer(peer string, ic IpnsCollector) error {
	cid, err := ic.GetIpnsCid(peer)
	row := []string{
		peer,
		time.Now().Format("2006-01-02 15:04:05"),
		cid,
	}
	err = AddRowToCSV(DATA_STORE_PATH_WEEK+WEEK_FILE, row)
	if err != nil {
		return err
	}

	return nil
}

func logIterations(iter int) error {
	row := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		strconv.Itoa(iter),
	}
	err := AddRowToCSV(DATA_STORE_PATH_WEEK+ITER_FILE, row)
	if err != nil {
		return err
	}
	return nil
}

func checkedAllOnes() error {
	row := []string{
		time.Now().Format("2006-01-02 15:04:05"),
	}
	err := AddRowToCSV("/home/axel/Desktop/done.csv", row)
	if err != nil {
		return err
	}
	return nil

}
