package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

type website struct {
	Name     string
	Worktree string
	Gitdir   string
	Token    string
}

type configuration struct {
	Website []website
}

var conf configuration

func init() {
	// read in config file
	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	blob := string(b)

	// parse toml
	if _, err := toml.Decode(blob, &conf); err != nil {
		log.Fatal(err)
	}

	// print out loaded websites
	for _, site := range conf.Website {
		fmt.Printf("loaded %s\n", site.Name)
	}

	// put websites into hashmap

}

func main() {
	http.HandleFunc("/", incomingHook)

	http.ListenAndServe(":8989", nil)
}

func incomingHook(rw http.ResponseWriter, r *http.Request) {
	// check if request is for a known repository

	// retrives the information for the repository

	// check if X-Gitlab-Token token matches

	// call procedure to update repo
}

func updateRepo(workDir string, repoDir string) {
	// set work and repo directory enviroment vars

	// set tmp work directory

	// Git pull and checkout -f and potentially LFS commands
	// N.B., that git must have a ssh key

	// run any post pull scripts in restricted shell or chroot

	// Do atomic mv to switch tmp work directory and live work dir
}
