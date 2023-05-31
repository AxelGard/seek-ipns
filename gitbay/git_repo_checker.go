package main

import (
	"context"
	"io/ioutil"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func isGitRepo(files []string) bool {
	for _, name := range files {
		//if name == ".git" || name == "README.md" {
		if strings.Contains(name, "README") {
			return true
		}
	}
	return false
}

func GetRepoFileData(cid string, sh shell.Shell, ctx context.Context) ([]byte, error) {
	f_list, err := sh.List(cid)
	if err != nil {
		return nil, err
	}
	for _, f := range f_list {
		cidf := cid + "/" + f.Name
		if cidf[0] != '/' {
			cidf = "/ipfs/" + cidf
		}
		fs, err := sh.FilesStat(ctx, cidf)
		if err != nil {
			return nil, err
		}
		if fs.Type == "directory" {
			sub_read, err := GetRepoFileData(cid+"/"+f.Name, sh, ctx)
			if err != nil {
				return nil, err
			}
			if sub_read != nil {
				return sub_read, nil
			}
		} else {
			if isGitRepo([]string{f.Name}) {
				r, err := sh.Cat(cid + "/" + f.Name)
				if err != nil {
					return nil, err
				}
				file_con, err := ioutil.ReadAll(r)
				if err != nil {
					return nil, err
				}
				return file_con, nil
			}
		}
	}
	return nil, nil
}
