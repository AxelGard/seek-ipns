package main

import (
	"fmt"
)

func main() {
	dir_cids := [4]string{
		"bafybeicnxhkmocvutxrexcwj62eqidgunz22wqmwzrrghtdyvi5vjgn6ci",
		"QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o", // NOT A DIR, as a test case
		"QmQmhYjzuJUzsM3uMVtByzsfdQG6H3LeGTkUUD1yHVf1vb",
		"QmPXeM8QwpqBcnzE54fduPb5mm9trMDfd2adgL1KrmNNP6",
	}
	for _, c := range dir_cids {
		files, err := GetFileNamesFromCID(c)
		if err != nil {
			panic(err)
		}
		if len(files) != 0 {
			fmt.Println(c)
			fmt.Println(files)
			fmt.Println(GetDataFromCID(c + "/" + files[0]))
			fmt.Println("----")
		}
	}
}
