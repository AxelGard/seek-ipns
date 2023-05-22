package main

func isGitRepo(files []string) bool {
	for _, name := range files {
		//if name == ".git" || name == "README.md" {
		if name == "README.md" {
			return true
		}
	}
	return false
}
