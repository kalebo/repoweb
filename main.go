package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
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
var websitemap map[string]*website
var updateHookRe = regexp.MustCompile("^/update/([a-zA-Z0-9-]*)/*$")

func init() {
	flag.Parse()

	// read in config file
	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		glog.Fatal(err)
	}
	blob := string(b)

	// parse toml
	if _, err := toml.Decode(blob, &conf); err != nil {
		glog.Fatal(err)
	}

	// put websites into hashmap for quick lookup
	websitemap = make(map[string]*website)
	for _, site := range conf.Website {
		websitemap[site.Name] = &site
		glog.Infof("Loaded %s\n", site.Name)
	}

}

func main() {
	http.HandleFunc("/update/", incomingHook)
	http.ListenAndServe(":8989", nil)
	fmt.Println(websitemap)
}

func incomingHook(w http.ResponseWriter, r *http.Request) {
	// check if request uri is valid
	captures := updateHookRe.FindStringSubmatch(r.URL.Path)
	if len(captures) != 2 {
		glog.Error("Invalid URI:", r.URL.Path)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	sitename := captures[1]

	// check if request is for a known repository
	if site, ok := websitemap[sitename]; ok {
		// retrives the information for the repository
		glog.Infof("Found site %s at %s", site.Name, site.Gitdir)

		// check if X-Gitlab-Token token matches
		if token := r.Header.Get("X-Gitlab-Token"); token != site.Token {
			glog.Infof("Token verification failed for %s", site.Name)
			http.Error(w, http.StatusText(403), 403)
			return
		}
		glog.Infof("Token verified for %s", site.Name)

		// call procedure to update repo
		if err := updateRepo(site.Gitdir, site.Worktree); err != nil {
			http.Error(w, fmt.Sprintf("Conflict: %s", err), 409)
			return
		}

	} else {
		http.Error(w, http.StatusText(404), 404)
		glog.Error("No such site:", sitename)
		return
	}
}

func updateRepo(gitDir string, workTree string) error {
	// set work and repo directory enviroment vars
	env := os.Environ()
	env = append(env, fmt.Sprintf("GIT_DIR=%s", gitDir))
	env = append(env, fmt.Sprintf("GIT_WORK_TREE=%s", workTree))

	// TODO: consider setting a tmp work directory

	// Git pull and checkout -f and potentially LFS commands
	// N.B., that git must have a ssh key
	statusCmd := exec.Command("git", "status", "-s")
	pullCmd := exec.Command("git", "fetch")
	checkoutCmd := exec.Command("git", "checkout", "-f", "master")

	statusCmd.Env = env
	pullCmd.Env = env
	checkoutCmd.Env = env

	// Check if there are unstaged changes or untracked files
	if output, _ := statusCmd.CombinedOutput(); len(output) > 1 {
		glog.Errorf("Work tree is not clean: %s", workTree)
		return errors.New("Work tree not clean")
	}

	// Fetch changes from up stream
	if output, err := pullCmd.CombinedOutput(); err != nil {
		glog.Errorf("Error running `git fetch` for %s: %s", gitDir, output)
		return errors.New("fetching repo failed")
	}

	// Checkout master
	if output, err := checkoutCmd.CombinedOutput(); err != nil {
		glog.Errorf("Error running `git checkout -f master` for %s in %s: %s", gitDir, workTree, output)
		return errors.New("checkout failed")
	}

	// TODO: run any post pull scripts in restricted shell or chroot?

	// TODO: consider doin an atomic mv to switch tmp work directory and live work dir
	return nil
}
