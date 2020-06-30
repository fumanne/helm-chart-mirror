package chart

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/mcuadros/go-version"
)

const (
	DESTINATION = "docs"
	MAXRETRY    = 5
)

var MIRROR_URL = fmt.Sprintf("https://%s.github.io/%s/", FetchGitEnv("GIT_USER", "fumanne"), FetchGitEnv("REPO_NAME", "helm-chart-mirror"))

type Maintainer struct {
	Email, Name string
}

type Chart struct {
	ApiVersion  string       `yaml:"apiVersion"`
	AppVersion  string       `yaml:"appVersion"`
	Created     time.Time    `yaml:"created"`
	Deprecated  bool         `yaml:"deprecated"`
	Description string       `yaml:"description"`
	Digest      string       `yaml:"digest"`
	Home        string       `yaml:"home"`
	Icon        string       `yaml:"icon"`
	Name        string       `yaml:"name"`
	Sources     []string     `yaml:"sources"`
	Urls        []string     `yaml:"urls"`
	Version     string       `yaml:"version"`
	Maintainers []Maintainer `yaml:"maintainers,omitempty"`
}

type Index struct {
	ApiVersion string              `yaml:"apiVersion"`
	Entries    map[string][]*Chart `yaml:"entries"`
	Generated  time.Time           `yaml:"generated"`
}

func (c *Chart) Download(wg *sync.WaitGroup, maxChan chan bool) {
	var u string
	path := prepare()
	defer wg.Done()
	defer func(maxChan chan bool) { <-maxChan }(maxChan)
	if len(c.Urls) > 1 {
		i := rand.Intn(len(c.Urls))
		u = c.Urls[i]
	} else {
		u = c.Urls[0]
	}

	//log.Printf("Download url %s", u)
	sliTarget := strings.Split(u, "/")
	target := filepath.Join(path, sliTarget[len(sliTarget)-1])

	resp, _ := _download(u, MAXRETRY)

	if err := resp.DownloadToFile(target); err != nil {
		log.Fatalf("Save File %s Failed, Error is %s\n", target, err)
	}
	c.setUrl(MIRROR_URL, sliTarget[len(sliTarget)-1])

}

func (c *Chart) setUrl(prefix, name string) {
	c.Urls = []string{prefix + name}
}

func prepare() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Get Current Path is Failed, Error is %s\n", err)
	}
	path := filepath.Join(dir, DESTINATION)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Create Dir %s Failed, Error is %s\n", path, err)
	}
	return path
}

func _download(url string, retry int) (*grequests.Response, error) {
	if retry == 0 {
		log.Printf("Download %s Retry Count is Max\n", url)
		return nil, errors.New("Max Count Reached!")
	}
	resp, err := grequests.Get(url, nil)
	if err != nil {
		log.Printf("Download  %s Failed, Error is %s\nTry Download it....Retry Count is %d \n", url, err, MAXRETRY-retry+1)
		time.Sleep(10 * time.Second)
		return _download(url, retry-1)
	}
	return resp, nil

}
