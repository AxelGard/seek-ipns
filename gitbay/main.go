package main

import (
	"fmt"
)

func main() {
	//SwarmCrawl()
	//RunCompare()
	//recollectData()
	//CrawlingEachNode()
	//CheckGitProvders()
	//test_GetFileNamesFromCid()

	ipns_rec := get_ipns_record_from_peer("12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi")
	fmt.Println(ipns_rec)
}
