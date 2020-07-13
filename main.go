package main

import (
	"gopkg.in/yaml.v2"
	"helm-chart-mirror/chart"
	"helm-chart-mirror/fetch"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"syscall"
)

const (
	OFFICAL_URL = "https://kubernetes-charts.storage.googleapis.com/index.yaml"
)

func main() {
	max := GetOpenFiles().Cur * 2 / 3
	log.Printf("Current Chan Size is %d\n", max)
	var wg sync.WaitGroup
	maxChan := make(chan bool, max)
	I := chart.Index{}
	if err := yaml.Unmarshal(fetch.FetchIndexYaml(OFFICAL_URL), &I); err != nil {
		log.Fatalf("Unmarshal struct %s Failed, Error is %s\n", I, err)
	}
	for _, charts := range I.Entries {
		for _, chart := range charts {
			wg.Add(1)
			maxChan <- true
			go chart.Download(&wg, maxChan)
		}
	}
	wg.Wait()
	I.SetGenerated()
	d, err := yaml.Marshal(&I)
	if err != nil {
		log.Fatalf("Marshal struct %s Failed, Error is %s\n", I, err)
	}
	if err := ioutil.WriteFile("./docs/index.yaml", d, os.ModePerm); err != nil {
		log.Fatalf("Write Index.yaml Failed, Error is %s\n", err)
	}

}

func GetOpenFiles() syscall.Rlimit {
	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		log.Fatalf("Get limit Error: %s\n", err)
	}
	log.Printf("Current Open Files Count is: %d\n", rlimit.Cur)
	log.Println("If you want to speed up download, Increase number by using `ulimit -n` command")
	return rlimit
}